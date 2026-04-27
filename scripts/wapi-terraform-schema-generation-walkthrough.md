# WAPI Schema-Driven Terraform Provider: Technical Walkthrough

**Status**: Proof of Concept — Demonstrated on `terraform-provider-nios`
**Implements**: [ADR: WAPI Schema Auto-Generation for Terraform](https://github.com/Infoblox-CTO/saas.docs/blob/main/docs/architecture/target-state/NIOS/WAPI/wapi-terraform-schema-generation.md)
**Provider Fork**: [acasella1984/terraform-provider-nios @ wapi-schema-plan-modifiers](https://github.com/acasella1984/terraform-provider-nios/tree/wapi-schema-plan-modifiers)

---

## The Problem

Terraform's `UseStateForUnknown` plan modifier tells the plan engine: _"this field's value hasn't changed — use whatever was in state last time."_ This suppresses noisy `(known after apply)` output. But applying it uniformly to all `Computed` fields is incorrect — fields derived from user inputs (like `dns_name` from `name`) produce the wrong plan value when the source field changes, causing:

```text
Error: Provider produced inconsistent result after apply
  .dns_name: was cty.StringVal("a-record2.example.com"),
  but now cty.StringVal("a-record4.example.com").
```

The root cause: the provider has no metadata to distinguish immutable fields (`creation_time` — never changes) from derived fields (`dns_name` — changes when `name` changes) from volatile fields (`last_queried` — changes on every DNS query).

## The Solution

A pipeline that extracts field-level mutability metadata from the NIOS source code and uses it to select the correct plan modifier per field:

| Field Classification | Modifier | Behavior |
| --- | --- | --- |
| Immutable after creation | `UseStateForUnknown` | Suppresses diff — value never changes |
| Client-derivable | `PunycodeDerivedFrom("name")` | Computes value at plan time from source field |
| Identity-encoding (ref) | `UseStateUnlessResourceChanges` | Preserves state unless other attributes changed |
| Server-recomputed on change | None | Shows `(known after apply)` — server determines value |
| Volatile | None | Shows `(known after apply)` — changes independently |
| Inherited (use_flag) | None | Parent value may change — can't predict |

## Results

| Metric | Value |
| --- | --- |
| WAPI object types extracted | 468 |
| Total fields classified | 5,240 |
| Provider model files modified | 406 |
| Computed fields with modifiers added | 2,895 |
| Custom plan modifiers created | 2 (PunycodeDerivedFrom, UseStateUnlessResourceChanges) |
| Build status | Clean (`go build ./...` passes) |
| Test status | All passing (`go test ./internal/... -short`) |
| Live validation | Create, update (name change), idempotency confirmed |

---

## 1. How Schema Extraction Works (`extract_schema.py`)

The extractor reads NIOS source code using Python's `ast` module — it does **not** import or execute any NIOS code. It cross-references three layers to produce metadata that the existing `?_schema` endpoint does not provide.

### The Three-Layer Pipeline

```text
Phase 1: RTXML (684 XML files)     → Database-level signals (key-type, unique-id, synthetic-direct)
Phase 2: pyabs (654 Python files)   → Links WAPI objects to RTXML structures via type mappings
Phase 3: WAPI  (260 Python files)   → Field names, types, CRUD permissions, i2w functions, use_flags
Phase 4: Classification             → Cross-references all three layers → mutability + computation
                                      Output: wapi_schema_all.json (468 types, 5,240 fields)
```

### Phase 1: RTXML — What the Database Knows

RTXML XML files define the OneDB database schema. Each `<member>` element has attributes that reveal mutability at the storage level:

```xml
<!-- products/dns/xml/bind_a.xml -->
<structure name="bind_a">
    <member name="address" type="rtxml.string" key-type="key" index="true"/>
    <member name="uuid" type="rtxml.string" unique-id="true"/>
    <member name="revision_id" type="rtxml.string" version-id="true"/>
    <member name="ddns_principal_lower" synthetic-direct="true"
            synthetic-direct-func="lowercase_ddns_principal"/>
</structure>
```

| XML Attribute | Mutability Signal | Example |
| --- | --- | --- |
| `key-type="key"` | **Immutable** — part of primary key, cannot change after creation | `address`, `name`, `zone` |
| `unique-id="true"` | **Immutable** — auto-generated UUID, set once | `uuid` |
| `version-id="true"` | **Volatile** — increments on every update | `revision_id` |
| `synthetic-direct="true"` | **Derived** — computed from other fields, not stored | `ddns_principal_lower` |
| `internal-ro="true"` | **Read-only** — server-managed, writes ignored | |

### Phase 2: pyabs — Linking WAPI to RTXML

Each WAPI object has a `cmdclass` that references a pyabs class, which has a `type` attribute mapping to the RTXML structure:

```text
WAPIObject_rra_2_14.cmdclass = a_record.ARecordCmd
    → ARecordCmd.objectcls = ARecord
        → ARecord.type = 'dns.bind_a'
            → RTXML structure ".com.infoblox.dns.bind_a"
```

The extractor follows this chain using AST parsing — it finds `type = 'dns.bind_a'` assignments in pyabs class definitions and builds the mapping without importing any modules.

### Phase 3: WAPI — What the API Knows

The WAPI layer defines the REST API contract using `WF*` field classes. The extractor parses these AST nodes to extract every keyword argument:

```python
# products/dns/server/src/wapi/rrobjs.py — WAPIObject_rra_2_14
fields = [
    WF('name', 'The FQDN...', iname='fqdn', create=True, search='=:~', std_field=True),
    WF('dns_name', '...', iname='dns_fqdn', purpose='r'),
    WFUInt('ttl', '...', use_flag='use_ttl', default={'value': None}),
    WFBool('disable', '...', iname='disabled', default={'value': False}),
    WFEnum('creator', '...', enumlist=['STATIC', 'DYNAMIC', 'SYSTEM']),
    WFTimeStamp('creation_time', '...', iname='creation_timestamp', purpose='r'),
    WFExtensibleAttributes(),
]
```

Each keyword argument maps to a schema property:

| WF Keyword | What It Tells Us | Schema Output |
| --- | --- | --- |
| `purpose='r'` | Read-only (no write/update) | `"supports": "r"`, `"mutability": "read_only"` |
| `create=True` | Required at creation | `"create_required": true` |
| `use_flag='use_ttl'` | Value inherited from parent unless flag overrides | `"computation": "inherited"` |
| `iname='dns_fqdn'` | Internal DB column name; `dns_fqdn` matches punycode pattern | `"computation": "client_derivable"` |
| `i2w=i2w_invert_boolean` | Converter function — checked against known transforms | `"derivation_function": "boolean_invert"` |
| `enumlist=[...]` | Allowed values for enum field | `"enum_values": [...]` |
| `create_only=True` | Settable only at creation | `"mutability": "immutable"` |

**Handling inherited fields**: Many objects inherit fields from abstract base classes via `fields.extend(BaseZone.base_fields)`. The parser does a two-pass approach — first collecting `base_fields` from abstract classes, then resolving `extend()` calls on concrete classes.

**Auto-generated fields**: The runtime creates `use_*` companion boolean fields for every field with `use_flag`. The parser replicates this by scanning for `use_flag` attributes and emitting the companion field.

### Phase 4: Cross-Layer Classification

The classifier applies signals in priority order to produce the final classification:

```text
1. RTXML: unique-id=true?       → read_only / server_computed / safe=true
2. RTXML: key-type=key?         → immutable / stored / safe=true
3. RTXML: synthetic-direct?     → read_only / derived / safe=false
4. WF: purpose='r' or 'rs'?    → read_only, then check:
   4a. iname in CLIENT_DERIVABLE_INAMES (dns_fqdn, dns_rdata)?
       → client_derivable / safe=true
   4b. i2w in CLIENT_DERIVABLE_TRANSFORMS (str.lower, invert_boolean)?
       → client_derivable / safe=true
   4c. wname in known immutable set (creation_time, creator, zone)?
       → server_computed / safe=true
   4d. Otherwise → server_computed / safe=false
5. WF: create_only=True?        → immutable / stored / safe=true
6. WF: use_flag set?            → inherited / safe=false
7. WF: has i2w function?        → derived / safe=false
8. Default                      → mutable / stored (user provides value)
```

**Output for `dns_name` on `record:a`**:

```json
{
    "name": "dns_name",
    "type": "string",
    "supports": "rs",
    "mutability": "read_only",
    "computation": "client_derivable",
    "derived_from": ["name"],
    "derivation_function": "punycode",
    "stable_after_create": true,
    "reconciliation_action": "compute_on_diff(punycode)"
}
```

Classification chain: No RTXML signal (WAPI-only field) → `purpose='r'` makes it read_only → `iname='dns_fqdn'` matches `CLIENT_DERIVABLE_INAMES` → `client_derivable` with `derivation_function=punycode`.

---

## 2. How Plan Modifier Recommendations Work (`gen_plan_modifiers.py`)

This script reads the extracted schema JSON and the provider's Go model files, then compares them to identify:

- **Missing modifiers**: Computed fields that need `UseStateForUnknown` or a custom modifier
- **Wrong modifiers**: Fields with `UseStateForUnknown` that should have a different (or no) modifier
- **Correct modifiers**: Fields where the current modifier matches the schema recommendation

### The Decision Matrix

The script maps each `computation` + `mutability` + `stable_after_create` combination to a Terraform plan modifier:

| `computation` | `mutability` | `safe` | Recommended Modifier |
| --- | --- | --- | --- |
| `stored` | `immutable` | `true` | `UseStateForUnknown` |
| `server_computed` | `read_only` | `true` | `UseStateForUnknown` |
| `server_computed` | `read_only` | `false` | None — value may change between applies |
| `client_derivable` | any | `true` | Custom modifier (`PunycodeDerivedFrom`, `LowercaseNormalized`) |
| `identity_encoding` | `read_only` | conditional | `UseStateUnlessResourceChanges` — preserves state unless other fields changed |
| `derived` | any | `false` | None — server recomputes using opaque logic |
| `inherited` | `mutable` | `false` | None — parent value may change |
| `stored` | `mutable` | N/A | Not `Computed` — user provides value |

### Client-Derivable Fields: Complete Inventory

Client-derivable fields follow standardized, deterministic transformations that the provider can replicate at plan time. These fields need **custom Go plan modifiers** instead of `UseStateForUnknown`. The schema extraction pipeline identifies them automatically via two detection methods:

1. **`CLIENT_DERIVABLE_INAMES`** — matches `iname` values like `dns_fqdn` to known punycode patterns
2. **`dns_*` prefix pattern** — any read-only field named `dns_X` where `X` is another field on the same object is classified as the punycode form of that field

The `apply_plan_modifiers.py` script handles these automatically: it adds `PunycodeDerivedFrom(source)` to fields without modifiers and replaces `UseStateForUnknown` with `PunycodeDerivedFrom(source)` on fields that have the wrong modifier.

#### Punycode Derivation — 37 fields in provider

The `dns_*` fields are the punycode/IDNA-encoded form of their source field, computed via RFC 5891.

**Go implementation**: `PunycodeDerivedFrom(sourceAttr)` — uses `golang.org/x/net/idna` to compute the value at plan time.

| Resource | Field | Source Field | Status |
| --- | --- | --- | --- |
| `record:a` | `dns_name` | `name` | Applied |
| `record:aaaa` | `dns_name` | `name` | Applied |
| `record:alias` | `dns_name` | `name` | Applied |
| `record:alias` | `dns_target_name` | `target_name` | Applied |
| `record:caa` | `dns_name` | `name` | Applied |
| `record:cname` | `dns_name` | `name` | Applied |
| `record:cname` | `dns_canonical` | `canonical` | Applied |
| `record:dname` | `dns_name` | `name` | Applied |
| `record:dname` | `dns_target` | `target` | Applied |
| `record:mx` | `dns_name` | `name` | Applied |
| `record:mx` | `dns_mail_exchanger` | `mail_exchanger` | Applied |
| `record:naptr` | `dns_name` | `name` | Applied |
| `record:naptr` | `dns_replacement` | `replacement` | Applied |
| `record:ns` | `dns_name` | `name` | Applied |
| `record:ptr` | `dns_name` | `name` | Applied |
| `record:ptr` | `dns_ptrdname` | `ptrdname` | Applied |
| `record:srv` | `dns_name` | `name` | Applied |
| `record:srv` | `dns_target` | `target` | Applied |
| `record:tlsa` | `dns_name` | `name` | Applied |
| `record:txt` | `dns_name` | `name` | Applied |
| `record:unknown` | `dns_name` | `name` | Applied |
| `ip_allocation` | `dns_name` | `name` | Applied |
| `sharedrecord:a` | `dns_name` | `name` | Applied |
| `sharedrecord:aaaa` | `dns_name` | `name` | Applied |
| `sharedrecord:cname` | `dns_name` | `name` | Applied |
| `sharedrecord:cname` | `dns_canonical` | `canonical` | Applied |
| `sharedrecord:mx` | `dns_name` | `name` | Applied |
| `sharedrecord:mx` | `dns_mail_exchanger` | `mail_exchanger` | Applied |
| `sharedrecord:srv` | `dns_name` | `name` | Applied |
| `sharedrecord:srv` | `dns_target` | `target` | Applied |
| `sharedrecord:txt` | `dns_name` | `name` | Applied |
| `dtc:record:cname` | `dns_canonical` | `canonical` | Applied |
| `zone_auth` | `dns_fqdn` | `fqdn` | Applied |
| `zone_auth` | `dns_soa_email` | `soa_email` | Applied |
| `zone_delegated` | `dns_fqdn` | `fqdn` | Applied |
| `zone_forward` | `dns_fqdn` | `fqdn` | Applied |
| `zone_stub` | `dns_fqdn` | `fqdn` | Applied |
| `zone_rp` | `dns_soa_email` | `soa_email` | Applied |

`UseStateForUnknown` has been removed from all `ref` fields (70 files) and replaced with the custom `UseStateUnlessResourceChanges` modifier.

#### Identity-Encoding Fields: `ref` — 70 resources

The `ref` field is the WAPI object reference — a base64-encoded path that includes the object's identity fields (name, FQDN, network, etc.). It changes whenever those identity fields change (e.g., renaming a record). This creates a unique problem:

- **`UseStateForUnknown`**: Causes "Provider produced inconsistent result after apply" when identity fields change — the plan preserves the old ref but WAPI returns a new one.
- **No modifier at all**: Causes spurious diffs on every plan — `ref` shows `(known after apply)` even when nothing changed, creating noise.

The solution is a custom plan modifier `UseStateUnlessResourceChanges` that compares all *known* planned attributes against state. If everything matches, the resource hasn't changed and `ref` is preserved. If any attribute differs, `ref` is left unknown so Terraform accepts the server's new value.

**Go implementation**: `internal/planmodifiers/ref/use_state_unless_changed.go`

```go
func (m useStateUnlessResourceChanges) PlanModifyString(
    _ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
    // On create — no prior state, leave unknown
    if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
        return
    }
    // If plan already has concrete value, don't override
    if !req.PlanValue.IsUnknown() {
        return
    }
    // Compare all known planned attributes against state
    planAttrs := map[string]tftypes.Value{}
    req.Plan.Raw.As(&planAttrs)
    stateAttrs := map[string]tftypes.Value{}
    req.State.Raw.As(&stateAttrs)

    for key, planVal := range planAttrs {
        if !planVal.IsKnown() { continue }  // Skip other Computed-only fields
        stateVal, exists := stateAttrs[key]
        if !exists || !planVal.Equal(stateVal) {
            return  // Something changed — leave ref unknown
        }
    }
    resp.PlanValue = req.StateValue  // Nothing changed — preserve state
}
```

**Behavior**:

| Scenario | Plan Shows | After Apply |
| --- | --- | --- |
| No changes | `ref = "record:a/..."` (preserved) | No diff |
| Name changed | `ref = (known after apply)` | New ref accepted |
| TTL changed | `ref = (known after apply)` | New ref accepted |
| Create | `ref = (known after apply)` | Server-generated ref |

Applied by `apply_plan_modifiers.py` via the `CONDITIONAL_STATE_FIELDS` set — every `ref` field on every resource gets this modifier automatically.

#### Boolean Inversion (`disable_discovery`) — 3 fields

The `disable_discovery` field is the boolean inverse of `enable_discovery`: `disable_discovery = !enable_discovery`.

**Go implementation**: `BooleanInvertedFrom(sourceAttr)` — created in `internal/planmodifiers/derived/boolean_invert.go`. Not yet applied to the provider model files because `disable_discovery` already has `Default: booldefault.StaticBool(false)` which handles the create case correctly.

| Resource | Field | Source Field | In Provider? | Status |
| --- | --- | --- | --- | --- |
| `fixedaddress` | `disable_discovery` | `enable_discovery` | Yes | Has Default (no modifier needed) |
| `ipv6fixedaddress` | `disable_discovery` | `enable_discovery` | Yes | Has Default (no modifier needed) |
| `record:host` | `disable_discovery` | `enable_discovery` | No | Not in provider |

#### Lowercase Normalization (`rrset_order`) — 1 field

| Resource | Field | Source Field | In Provider? | Status |
| --- | --- | --- | --- | --- |
| `record:host` | `rrset_order` | `rrset_order` (self) | No | Not in provider |

### Audit Mode

```bash
python3 gen_plan_modifiers.py --audit --schema-dir output \
    --provider-dir /path/to/terraform-provider-nios --objects record:a
```

Output:

```text
record:a — AUDIT: Current vs Recommended Plan Modifiers
  WRONG MODIFIERS:
    ✗ dns_name: has UseStateForUnknown, needs PunycodeDerivedFrom("name")
  MISSING MODIFIERS:
    ⚠ cloud_info: needs UseStateForUnknown
    ⚠ ddns_principal: needs UseStateForUnknown
  CORRECT:
    ✓ creation_time: UseStateForUnknown (safe — immutable after creation)
    ✓ zone: UseStateForUnknown (safe — immutable after creation)
    ✓ ref: no modifier (correct — always recomputed)
```

---

## 3. How Plan Modifiers Are Applied (`apply_plan_modifiers.py`)

This script scans all `model_*.go` files under `internal/service/` and injects the correct plan modifier for every `Computed: true` field.

### Which Files Are Modified and Why

The terraform-provider-nios has ~520 `model_*.go` files. Each file defines a Go struct (the data model) and a `map[string]schema.Attribute` variable (the Terraform schema). For example, `model_record_a.go`:

```go
// The schema definition that Terraform reads:
var RecordAResourceSchemaAttributes = map[string]schema.Attribute{
    "creation_time": schema.Int64Attribute{
        Computed:            true,
        MarkdownDescription: "The time of the record creation.",
        // ← THIS IS WHERE the script injects PlanModifiers
    },
}
```

The script modifies **only the schema attribute blocks** within these files. It never touches:

- Model struct definitions
- `Expand()` functions (Terraform model → API request)
- `Flatten()` functions (API response → Terraform model)
- Resource CRUD functions (Create, Read, Update, Delete)
- Validators, Defaults, or any other existing schema properties

### How the Script Finds and Modifies Fields

For each `model_*.go` file, the script:

**1. Scans for `Computed: true` lines:**

```go
"creation_time": schema.Int64Attribute{    // ← detects Int64 attribute type
    Computed:            true,              // ← found Computed: true
    MarkdownDescription: "...",
},                                          // ← finds closing brace
```

**2. Checks skip conditions:**

- Already has `PlanModifiers:` block → skip (don't double-apply)
- Has `Default:` → skip (Computed+Optional+Default fields are self-managing)
- Field name is in `ALWAYS_RECOMPUTED_FIELDS` (`id`) → skip
- Field name is in `CONDITIONAL_STATE_FIELDS` (`ref`) → apply custom `UseStateUnlessResourceChanges` modifier

**3. Consults the WAPI schema for safety:**

The script guesses the WAPI type from the filename (`model_record_a.go` → `record:a`) and looks up the field:

```python
# If schema says this field is unsafe, skip it
if safety_info.get("stable_after_create") is False:
    continue  # Don't add UseStateForUnknown to volatile/derived fields
```

**4. Injects the correct typed modifier** based on the attribute type:

| Go Attribute Type | Import Added | Modifier Injected |
| --- | --- | --- |
| `schema.StringAttribute` | `stringplanmodifier` | `stringplanmodifier.UseStateForUnknown()` |
| `schema.Int64Attribute` | `int64planmodifier` | `int64planmodifier.UseStateForUnknown()` |
| `schema.BoolAttribute` | `boolplanmodifier` | `boolplanmodifier.UseStateForUnknown()` |
| `schema.MapAttribute` | `mapplanmodifier` | `mapplanmodifier.UseStateForUnknown()` |
| `schema.SingleNestedAttribute` | `objectplanmodifier` | `objectplanmodifier.UseStateForUnknown()` |
| `schema.ListAttribute` | `listplanmodifier` | `listplanmodifier.UseStateForUnknown()` |

**5. Adds import statements** to the file's `import` block for any new planmodifier packages.

### Concrete Example: What Changes in `model_record_a.go`

Before (field has no plan modifier):

```go
"cloud_info": schema.SingleNestedAttribute{
    Attributes:          RecordACloudInfoResourceSchemaAttributes,
    Computed:            true,
    MarkdownDescription: "The cloud information associated with the record.",
},
```

After (script injects modifier + import):

```go
"cloud_info": schema.SingleNestedAttribute{
    Attributes:          RecordACloudInfoResourceSchemaAttributes,
    Computed:            true,
    MarkdownDescription: "The cloud information associated with the record.",
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.UseStateForUnknown(),
    },
},
```

And at the top of the file, the import block gains:

```go
import (
    // ... existing imports ...
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
)
```

### Fields That Require Custom Plan Modifiers

The bulk script handles `UseStateForUnknown` for ~2,890 fields and `UseStateUnlessResourceChanges` for `ref` on 70 resources. The following field types require **custom derivation** plan modifiers:

| Field | Resource(s) | Modifier | Why |
| --- | --- | --- | --- |
| `dns_name` | All record types | `PunycodeDerivedFrom("name")` | Derived from `name` via RFC 5891 punycode |
| `dns_canonical` | `record:cname` | `PunycodeDerivedFrom("canonical")` | Same pattern, different source field |

These are identified by `gen_plan_modifiers.py --audit` as "WRONG MODIFIERS" and are handled automatically by `apply_plan_modifiers.py`.

---

## 4. The dns_name Inconsistency Bug: Reproduction and Fix

This section documents the exact bug from the NIOS Plugin Discussion PDF and how the schema-driven approach fixes it.

### The Bug

When a user changes `name` on a `record:a` (e.g., `a-record2` → `a-record4`):

| Field | Old State | `UseStateForUnknown` Plan | Server Returns | Match? |
| --- | --- | --- | --- | --- |
| `name` | `a-record2.example.com` | `a-record4.example.com` | `a-record4.example.com` | ✅ |
| `dns_name` | `a-record2.example.com` | `a-record2.example.com` | `a-record4.example.com` | ❌ |
| `ref` | `record:a/...a-record2...` | `record:a/...a-record2...` | `record:a/...a-record4...` | ❌ |

Terraform compares planned values against actual values after apply. Mismatch → `Provider produced inconsistent result after apply`.

### The Fix

Two changes to `model_record_a.go`:

```go
// dns_name: Replace UseStateForUnknown with PunycodeDerivedFrom
"dns_name": schema.StringAttribute{
    Computed: true,
    PlanModifiers: []planmodifier.String{
        derivedmod.PunycodeDerivedFrom("name"),   // computes correct value at plan time
    },
},

// ref: Replace UseStateForUnknown with UseStateUnlessResourceChanges
"ref": schema.StringAttribute{
    Computed: true,
    PlanModifiers: []planmodifier.String{
        refmod.UseStateUnlessResourceChanges(),    // preserves state unless other fields changed
    },
},
```

### Plan Output After Fix

When identity fields change (rename):
```text
  ~ resource "nios_dns_record_a" "test" {
      ~ dns_name = "a-record4.example.com" -> "a-record5.example.com"   # correct!
      ~ name     = "a-record4.example.com" -> "a-record5.example.com"
      ~ ref      = "record:a/..." -> (known after apply)                 # correct!
        # (12 unchanged attributes hidden)                                # suppressed!
    }
```

When nothing changes (idempotent re-plan):
```text
No changes. Your infrastructure matches the configuration.
```

The `UseStateUnlessResourceChanges` modifier ensures `ref` is preserved when no other attributes changed, eliminating spurious diffs while still allowing Terraform to accept a new ref value when identity fields change.

---

## 5. All Changes Applied to Test Resources

### Test Configuration

The test `main.tf` (in `examples/`) exercises four resources covering all classification categories:

- `nios_dns_zone_auth` — immutable key fields, server-computed fields
- `nios_dns_record_a` — client-derivable `dns_name`, volatile `ref`
- `nios_dns_record_a` (with TTL) — inherited fields via `use_flag`
- `nios_dns_record_cname` — client-derivable `dns_name` and `dns_canonical`

### Changes per File

**`model_record_a.go`** — 3 fields changed:

| Field | Before | After | Schema Signal |
| --- | --- | --- | --- |
| `dns_name` | `UseStateForUnknown` | `PunycodeDerivedFrom("name")` | `computation=client_derivable, iname=dns_fqdn` |
| `cloud_info` | No modifier | `UseStateForUnknown` | `computation=server_computed, mutability=read_only` |
| `ddns_principal` | No modifier | `UseStateForUnknown` | `computation=stored, Optional+Computed` |

**`model_record_cname.go`** — 6 fields changed:

| Field | Before | After | Schema Signal |
| --- | --- | --- | --- |
| `dns_name` | `UseStateForUnknown` | `PunycodeDerivedFrom("name")` | `computation=client_derivable` |
| `dns_canonical` | `UseStateForUnknown` | `PunycodeDerivedFrom("canonical")` | `computation=client_derivable, iname=dns_fqdn` |
| `ref` | `UseStateForUnknown` | `UseStateUnlessResourceChanges` | Always recomputed — `CONDITIONAL_STATE_FIELDS` |
| `cloud_info` | No modifier | `UseStateForUnknown` | `computation=server_computed` |
| `ddns_principal` | No modifier | `UseStateForUnknown` | `computation=stored, Optional+Computed` |
| `extattrs_all` | `AssociateInternalId` only | + `UseStateForUnknown` | Consistency with record:a pattern |

**`apply_plan_modifiers.py`** — 3 tooling updates:

| Change | Why |
| --- | --- |
| Added `CONDITIONAL_STATE_FIELDS` set (`ref`) with `UseStateUnlessResourceChanges` modifier | Preserves ref in state when no other attributes changed; leaves unknown on identity changes |
| `ALWAYS_RECOMPUTED_FIELDS` now contains only `id` | `ref` moved to `CONDITIONAL_STATE_FIELDS` for smarter handling |
| Fixed `build_field_safety_map` to handle `extract_schema.py` output format | Schema lookup was silently failing — all fields got UseStateForUnknown as fallback |

### Remaining `(known after apply)` on Create

After all changes, these fields still show `(known after apply)` on **create** operations. This is correct — they are genuinely server-determined:

| Field | Why Unpredictable |
| --- | --- |
| `ref` | Server-generated base64-encoded path |
| `creation_time` | Server timestamp at moment of creation |
| `zone` | Depends on which zones exist on the appliance |
| `extattrs_all` | EA store populates after create (includes auto-generated Terraform Internal ID) |
| `last_queried` | Server-only counter, changes on every DNS query |
| `cloud_info`, `aws_rte53_record_info`, `discovered_data`, `ms_ad_user_data` | Complex structs from optional integrations |
| `reclaimable`, `shared_record_group`, `ddns_principal` | Server determines based on configuration |
| `ttl` (when `use_ttl=false`) | Inherited from zone/grid hierarchy |

`UseStateForUnknown` only helps on **update** plans (where prior state exists). On **create**, there is no prior state. This is standard Terraform behavior — AWS, Azure, and GCP providers all show `(known after apply)` for `arn`, `id`, `created_at`, etc.

---

## 6. Extraction Accuracy

The accuracy percentages measure how well the static AST parser reproduces what a running NIOS instance reports via the `?_schema` endpoint.

- **Name Match**: _Of all field names the running NIOS reports, what percentage did the extractor find?_ 100% means every field was discovered.
- **Field Match**: _Of common fields, what percentage have identical attributes (type, supports, searchable_by)?_ Mismatches are primarily cosmetic — type labeling (`struct` vs specific wapitype name) and domain-specific WF class overrides.

| Object | Name Match | Field Match | Gap Cause |
| --- | --- | --- | --- |
| `record:a` | 100% | 83% | Type naming (struct vs awsrte53recordinfo) |
| `record:cname` | 100% | 91% | Type naming |
| `zone_auth` | 94% | 75% | Missing WFFuncCall endpoints + type naming |
| `zone_forward` | 97% | 90% | Missing 1 FuncCall endpoint |
| `zone_delegated` | 96% | 93% | Missing 1 FuncCall endpoint |
| `network` | 93% | 81% | Missing FuncCall endpoints + DHCP module extends |
| `fixedaddress` | 100% | 83% | Type naming |
| `view` | 100% | 86% | Type naming |

A field match of 83% does **not** mean 17% of fields produce wrong Terraform behavior. The mismatches are type labels (`struct` vs `awsrte53recordinfo`) and minor `supports` flag differences. The core metadata driving plan modifier selection (`mutability`, `computation`, `stable_after_create`) is correct for essentially all fields because it uses multiple cross-layer signals, not any single attribute.

### Field Classification Summary (468 types, 5,240 fields)

| Classification | Count | Terraform Impact |
| --- | --- | --- |
| mutable/stored | 2,839 | `Optional` or `Required` — user provides value |
| read_only/server_computed | 1,532 | `Computed` — UseStateForUnknown only if safe |
| mutable/inherited | 667 | `Optional` — inherits from parent if use_flag=false |
| immutable/stored | 225 | `Required` + `ForceNew` (RequiresReplace) |
| derived | 172 | `Computed` — NO UseStateForUnknown |
| client_derivable | 30 | `Computed` + custom plan modifier |
| UseStateForUnknown safe | 502 | Can safely suppress plan diffs |
| UseStateForUnknown unsafe | 4,738 | Must NOT suppress |

---

## 7. WAPI Changes Needed for Native Schema Support

The static AST extraction works today but has inherent limitations. Extending the WAPI `?_schema` endpoint to expose classification metadata natively would make the pipeline 100% accurate and eliminate the need for NIOS source access.

### What's Missing from `?_schema` v2

The current endpoint returns:

```json
{"name": "dns_name", "type": ["string"], "supports": "r"}
```

It does not include: `create_required`, `create_only`, `nillable`, `use_flag`, `default_value`, `iname`, `mutability`, `computation`, `derived_from`, `derivation_function`, `stable_after_create`, or `i2w_function`.

### Proposed `?_schema` v3

Extend `WapiFieldInfo_2._produce_field_info()` in `wapibaseobj.py` to emit the additional metadata. This requires ~30 lines of Python in a new `WapiFieldInfo_3` subclass plus the classification function. The runtime has direct access to live WF field objects (not AST representations), making classification 100% accurate — including fields modified by `modify_fields()` at runtime, which the static parser cannot trace.

| Aspect | Static AST Extraction | Native `?_schema` v3 |
| --- | --- | --- |
| Accuracy | ~83-93% field match | 100% — runs inside the runtime |
| Source access | Requires NIOS source checkout | Only needs a running instance |
| `modify_fields()` | Cannot execute (AST limitation) | Already executed by metaclass |
| Maintenance | Script must track code pattern changes | Built into NIOS — updates automatically |

See the [ADR](https://github.com/Infoblox-CTO/saas.docs/blob/main/docs/architecture/target-state/NIOS/WAPI/wapi-terraform-schema-generation.md) for the full `?_schema` v3 specification and implementation details.

---

## 8. SDK-Level ToMap() Fixes (`apply_sdk_omitempty.py`)

During live testing, the schema-driven pipeline uncovered a class of bugs where the Go SDK serializes zero/empty values that WAPI rejects. The root cause is in the SDK's auto-generated `ToMap()` methods, which use only `!IsNil()` guards — insufficient for struct pointers, slices, maps, and enum/object_ref strings.

### The Problem

The Go SDK (`infoblox-nios-go-client`) generates `ToMap()` methods that serialize every non-nil field:

```go
// Generated SDK code — sends zero-value structs and empty strings
func (o Network) ToMap() (map[string]interface{}, error) {
    toSerialize := map[string]interface{}{}
    if !IsNil(o.MsAdUserData) {
        toSerialize["ms_ad_user_data"] = o.MsAdUserData  // sends {} for zero struct
    }
    if !IsNil(o.KnownClientsOption) {
        toSerialize["known_clients_option"] = o.KnownClientsOption  // sends "" for empty enum
    }
    if !IsNil(o.DiscoveryMember) {
        toSerialize["discovery_member"] = o.DiscoveryMember  // sends "" for empty object ref
    }
    return toSerialize, nil
}
```

This causes WAPI errors:

| Field Type | Error | Root Cause |
| --- | --- | --- |
| Struct pointer (`*NetworkMsAdUserData`) | `Field is not writable: ms_ad_user_data` | Zero-value struct `{}` sent |
| Enum string (`*string`) | `Invalid enum value for known_clients_option` | Empty string `""` sent |
| Object ref string (`*string`) | `Grid Member  not found` (double space = empty) | Empty string `""` sent |
| Slice (`[]MonitorMember`) | Various | Empty `[]` sent |
| Map (`map[string]interface{}`) | Various | Empty `{}` sent |

### The Fix: `apply_sdk_omitempty.py`

This script modifies `ToMap()` methods across the entire SDK to add zero-value checks. It operates on the vendored SDK at `vendor/github.com/infobloxopen/infoblox-nios-go-client/`:

| Field Type | Before | After |
| --- | --- | --- |
| Struct pointer | `if !IsNil(o.X) {` | `if !IsNil(o.X) { m, _ := o.X.ToMap(); if len(m) > 0 {` |
| Slice | `if !IsNil(o.X) {` | `if !IsNil(o.X) && len(o.X) > 0 {` |
| Map pointer | `if !IsNil(o.X) {` | `if !IsNil(o.X) && len(*o.X) > 0 {` |
| Enum string | `if !IsNil(o.X) {` | `if !IsNil(o.X) && *o.X != "" {` |
| Object ref string | `if !IsNil(o.X) {` | `if !IsNil(o.X) && *o.X != "" {` |
| Regular string | `if !IsNil(o.X) {` | Unchanged — `""` is valid (clears field value) |

**Key design decision**: Regular string fields (like `comment`) are left unchanged because sending `""` to WAPI is how you clear a field's value. Only enum and object_ref strings get the empty-string guard, because `""` is never a valid enum value or object reference.

#### Object Reference Detection via `WFDBObjectByString`

Some fields are typed as `string` in the Go SDK but actually reference WAPI objects by name (grid members, views, policies). These are defined using `WFDBObjectByString` in the WAPI source. The script promotes these to `"object_ref"` type so they get the `!= ""` guard:

```python
OBJECT_REF_WF_CLASSES = frozenset({"WFDBObjectByString"})
# If schema says wf_class is WFDBObjectByString, treat as object_ref
if ftype == "string" and wf_class in OBJECT_REF_WF_CLASSES:
    ftype = "object_ref"
```

Examples: `discovery_member`, `client_cert`, `failover_association`, `record_name_policy`.

### Results

| Metric | Value |
| --- | --- |
| SDK model files scanned | ~1,503 |
| ToMap() methods modified | ~1,500 |
| Struct pointer fields fixed | ~3,200 |
| Slice fields fixed | ~800 |
| Enum/object_ref string fields fixed | ~400 |
| Errors eliminated | All zero-value serialization errors |

### Usage

```bash
# Apply SDK ToMap() fixes to vendored SDK
python3 apply_sdk_omitempty.py \
    --sdk-dir /path/to/terraform-provider-nios/vendor/github.com/infobloxopen/infoblox-nios-go-client \
    --schema-dir /path/to/schema_output \
    --dry-run

# Apply for real
python3 apply_sdk_omitempty.py \
    --sdk-dir /path/to/terraform-provider-nios/vendor/github.com/infobloxopen/infoblox-nios-go-client \
    --schema-dir /path/to/schema_output
```

---

## 9. Provider Expand-Level Fixes (Historical)

> **Note**: This section documents the original provider-level workaround. For new work, use `apply_sdk_omitempty.py` (Section 8) instead.

During initial testing, read-only fields were commented out of `Expand` functions with `TODO(SDK)` markers. These fixes are **no longer needed** when `apply_sdk_omitempty.py` is applied to the vendored SDK, because the SDK's `ToMap()` correctly skips zero-value fields.

### What Was Found

The WAPI schema classifies fields like `ms_ad_user_data`, `rir_organization`, `cloud_info`, and `discovered_data` as `supports='r'` (read-only). The provider's `Expand` function includes these fields in every create/update call. When the Go SDK serializes the struct, empty/nil values for these read-only fields get sent to WAPI, which rejects them:

| Field | Error | Root Cause |
| --- | --- | --- |
| `ms_ad_user_data` | `Field is not writable: ms_ad_user_data` | Read-only struct sent as `{}` in update |
| `rir_organization` | `RIR organization object  not found` | Empty string sent for read-only ref field |
| `discovery_member` | `Grid Member  not found` | Empty WFDBObjectByString sent as `""` |
| `known_clients` | `Invalid enum value for known_clients_option` | Empty string sent for enum field |
| `content_check_op` | `Invalid value for content_check_op` | Empty string sent for enum field |
| `topology` (on lbdn) | `Invalid reference` | Empty string sent for read-only ref |

### How the Schema Pipeline Detected These

The `WF_CLASS_FORCED_PURPOSE` fix in `extract_schema.py` correctly classified `WFMSADUserData`, `WFCloudInfo`, `WFDiscoveryData`, and `WFAWSRte53RecordInfo` as `purpose='r'`. This made the schema output `supports='r'` / `mutability='read_only'` for 79 fields across all object types.

### Important Distinction: User-Specified vs Server-Computed Fields

**Rule**: Only exclude fields from Expand that are **server-computed read-only** (`supports='r'`). Fields that the user specifies in their `.tf` config (like `servers`, `member`, `pools`) must remain in Expand even if they cause issues — the plan/apply contract requires the result to match what the user declared.

| Field Type | In Expand? | Example |
| --- | --- | --- |
| Server-computed read-only (`supports='r'`) | **Exclude** | `ms_ad_user_data`, `rir_organization` |
| User-specified writable (`supports='rwu'`) | **Keep** | `servers`, `pools`, `member` |
| Optional writable with empty-string bug | **Keep, fix in SDK** | `known_clients`, `content_check_op` |

### Long-Term SDK Fix Options

The Go SDK (`infoblox-nios-go-client`) currently generates flat structs where every field is serialized regardless of write permissions:

#### Option A: Separate Read/Write model structs

Generate `NetworkWritable` (only `supports` containing `w` or `u`) and `NetworkReadable` (all fields) from the WAPI schema JSON. The SDK's `Create()` and `Update()` methods use the writable struct; `Read()` uses the readable struct.

#### Option B: Custom JSON marshaler with field filter

Add a `MarshalForWrite()` method that consults a compile-time map of writable fields, generated from the schema.

#### Option C: `wapi:"readonly"` struct tag

Add a custom struct tag to read-only fields. The SDK's HTTP client strips tagged fields before serialization.

Option A aligns with the CBA report's call for _"systematic use of Go SDK autogeneration from the WAPI API"_ and the same `wapi_schema_all.json` that drives plan modifier generation can drive SDK struct generation.

### Remaining SDK-Level Issues

| Issue | Symptom | SDK Fix Needed |
| --- | --- | --- |
| `bfdtemplate` deprecated field | `authentication_type` deprecated in NIOS 9.1.0 but SDK still sends it | SDK should support version-aware field exclusion from schema metadata |
| Server default values not declared | `ratio` returns `1` when unset, plan/apply inconsistency | Fields with `has_server_default=true` in schema need `Optional: true, Computed: true` |

### Server-Default Fields and `has_server_default`

The schema extraction pipeline now emits `has_server_default: true` for fields where the WAPI server provides a default value (e.g., `ratio` defaults to `1`, `ttl` defaults to zone TTL). When these fields aren't set by the user, Terraform sees `null` in the plan but the server returns the default value, causing "Provider produced inconsistent result after apply."

The fix is straightforward: fields with `has_server_default=true` must be declared as `Optional: true, Computed: true` in the Terraform schema, with a `UseStateForUnknown` plan modifier. The schema JSON contains 2,220+ fields with `has_server_default=true`, of which ~1,786 are writable.

---

## 10. Schema Field Definitions

The extracted WAPI schema produces the following metadata for each field. These are the fields in the `wapi_schema_all.json` output:

| Schema Field | Type | Description |
| --- | --- | --- |
| `name` | string | WAPI field name (e.g., `dns_name`, `comment`, `ttl`) |
| `type` | string | Field data type: `string`, `uint`, `int`, `bool`, `enum`, `struct`, `timestamp`, `object_ref`, `list`, `map` |
| `wf_class` | string | WAPI WF class name (e.g., `WF`, `WFUInt`, `WFBool`, `WFEnum`, `WFTimeStamp`, `WFDBObjectByString`) |
| `supports` | string | CRUD permission string: `r` (read), `w` (write/create), `u` (update), `s` (search). E.g., `rwu`, `rs`, `r` |
| `iname` | string | Internal database column name mapping (e.g., `dns_fqdn`, `disabled`, `creation_timestamp`) |
| `create_required` | bool | Whether the field is required at object creation (`create=True` in WF definition) |
| `create_only` | bool | Field is settable only at creation, immutable after (`create_only=True`) |
| `nillable` | bool | Whether the field can be set to null/empty to clear its value |
| `use_flag` | string | Name of companion boolean field that controls inheritance (e.g., `use_ttl` for `ttl`) |
| `default_value` | any | Default value assigned by the server when unset |
| `has_server_default` | bool | Whether WAPI provides a default value for this field |
| `enum_values` | list | Allowed values for enum-type fields (e.g., `["STATIC", "DYNAMIC", "SYSTEM"]`) |
| `searchable_by` | string | Search operators supported (e.g., `=`, `~`, `:=`, `:~`) |
| `standard_field` | bool | Whether the field is returned by default without `_return_fields` |
| `mutability` | string | Classification: `mutable`, `immutable`, `read_only` |
| `computation` | string | How the value is determined: `stored`, `server_computed`, `client_derivable`, `derived`, `inherited` |
| `derived_from` | list | Source field(s) for derived/client_derivable fields (e.g., `["name"]` for `dns_name`) |
| `derivation_function` | string | Transform applied: `punycode`, `boolean_invert`, `lowercase` |
| `stable_after_create` | bool | Whether the value is safe to cache in state (true = `UseStateForUnknown` safe) |
| `reconciliation_action` | string | Terraform plan modifier recommendation (e.g., `use_state_for_unknown`, `compute_on_diff(punycode)`) |
| `i2w_function` | string | Internal-to-WAPI converter function name from NIOS source |
| `rtxml_key_type` | string | RTXML database key classification: `key`, `pkey`, `natural-key` |
| `rtxml_unique_id` | bool | Whether the field is a unique identifier in the database |
| `rtxml_synthetic_direct` | bool | Whether the field is computed from other fields at the database level |

---

## 11. Reproducing This Workflow

```bash
# 0. Fix update
python3 add_put_expand.py --repo-path /Users/vchandra/wspace/terraform-provider-nios

# 1. Extract schema
python3 extract_schema.py --nios-root /path/to/nios --output-dir output

python extract_schema.py --nios-root ./nios --output-dir ./output \
        --objects record:a

# 2. Validate extraction accuracy (optional)
python3 validate_schema.py --schema-dir output/wapi_schema_all.json \
    --diff-v2-schema /path/to/nios/.../wapischemas/v2/wapi_v2.13.8.schemav2

python3 validate_schema.py --schema-dir output/wapi_schema_all.json \
    --diff-v2-schema ./nios/products/tests/server/src/bin/harness/datasets/wapischemas/v2/wapi_v2.13.7.schemav2  

# 3. Audit current provider state
python3 gen_plan_modifiers.py --audit --schema-dir output \
    --provider-dir /path/to/terraform-provider-nios

python3 gen_plan_modifiers.py --audit --schema-dir output \
    --provider-dir /Users/vchandra/wspace/terraform-provider-nios > audit.txt

# 4. Apply SDK ToMap() fixes (fixes zero-value serialization bugs)
python3 apply_sdk_omitempty.py \
    --sdk-dir /path/to/provider/vendor/github.com/infobloxopen/infoblox-nios-go-client \
    --schema-dir output

python3 apply_sdk_omitempty.py \
    --sdk-dir /Users/vchandra/wspace/terraform-provider-nios/vendor/github.com/infobloxopen/infoblox-nios-go-client \
    --schema-dir output

# 5. Apply bulk plan modifiers (dry-run first)
python3 apply_plan_modifiers.py --provider-dir /path/to/provider \
    --schema-dir output --dry-run
python3 apply_plan_modifiers.py --provider-dir /path/to/provider \
    --schema-dir output

python3 apply_plan_modifiers.py --provider-dir /Users/vchandra/wspace/terraform-provider-nios --schema-dir output

# 6. The scripts now handle all modifier types automatically:
#    - UseStateForUnknown for safe immutable/server-computed fields
#    - PunycodeDerivedFrom for dns_* client-derivable fields
#    - UseStateUnlessResourceChanges for ref fields
#    - Skips id fields (ALWAYS_RECOMPUTED_FIELDS)
#    - SDK ToMap() fixes for struct/slice/map/enum/object_ref zero values

# 7. Build, test, validate
cd /path/to/terraform-provider-nios
GOTOOLCHAIN=auto go mod vendor && go build ./... && go test ./internal/... -short
terraform plan

# 8. (For version upgrades) Diff schemas between NIOS versions
python3 validate_schema.py --schema-dir output \
    --version-diff --old-schema output-9.0.6/wapi_schema_all.json \
    --new-schema output-9.1.0/wapi_schema_all.json

# 9. (For CI/CD) Validate plan output
terraform plan -json > plan.json
python3 validate_schema.py --schema-dir output/wapi_schema_all.json \
    --plan-snapshot plan.json
```

---

## 12. Pipeline Tool Reference

| Tool | Purpose | Replaces |
| --- | --- | --- |
| `extract_schema.py` | AST-based WAPI schema extraction from NIOS source | Manual schema reading |
| `validate_schema.py` | Diff schema against v2, live WAPI, Terraform provider, or between versions | Manual auditing |
| `gen_plan_modifiers.py` | Audit/recommend plan modifiers from schema classification | Manual modifier selection |
| `apply_plan_modifiers.py` | Inject correct plan modifiers into provider model files | Manual Go editing |
| `apply_sdk_omitempty.py` | Fix SDK ToMap() zero-value serialization | `apply_expand_fixes.py` (deprecated) |
| ~~`apply_expand_fixes.py`~~ | ~~Comment out read-only fields in Expand functions~~ | **Deprecated** — use `apply_sdk_omitempty.py` |
