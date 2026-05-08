package dtc

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dtc"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	internaltypes "github.com/infobloxopen/terraform-provider-nios/internal/types"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type DtcLbdnModel struct {
	Ref                      types.String                     `tfsdk:"ref"`
	AuthZones                internaltypes.UnorderedListValue `tfsdk:"auth_zones"`
	AutoConsolidatedMonitors types.Bool                       `tfsdk:"auto_consolidated_monitors"`
	Comment                  types.String                     `tfsdk:"comment"`
	Disable                  types.Bool                       `tfsdk:"disable"`
	ExtAttrs                 types.Map                        `tfsdk:"extattrs"`
	ExtAttrsAll              types.Map                        `tfsdk:"extattrs_all"`
	Health                   types.Object                     `tfsdk:"health"`
	LbMethod                 types.String                     `tfsdk:"lb_method"`
	Name                     types.String                     `tfsdk:"name"`
	Patterns                 internaltypes.UnorderedListValue `tfsdk:"patterns"`
	Persistence              types.Int64                      `tfsdk:"persistence"`
	Pools                    types.List                       `tfsdk:"pools"`
	Priority                 types.Int64                      `tfsdk:"priority"`
	Topology                 types.String                     `tfsdk:"topology"`
	Ttl                      types.Int64                      `tfsdk:"ttl"`
	Types                    internaltypes.UnorderedListValue `tfsdk:"types"`
	UseTtl                   types.Bool                       `tfsdk:"use_ttl"`
}

var DtcLbdnAttrTypes = map[string]attr.Type{
	"ref":                        types.StringType,
	"auth_zones":                 internaltypes.UnorderedListOfStringType,
	"auto_consolidated_monitors": types.BoolType,
	"comment":                    types.StringType,
	"disable":                    types.BoolType,
	"extattrs":                   types.MapType{ElemType: types.StringType},
	"extattrs_all":               types.MapType{ElemType: types.StringType},
	"health":                     types.ObjectType{AttrTypes: DtcLbdnHealthAttrTypes},
	"lb_method":                  types.StringType,
	"name":                       types.StringType,
	"patterns":                   internaltypes.UnorderedListOfStringType,
	"persistence":                types.Int64Type,
	"pools":                      types.ListType{ElemType: types.ObjectType{AttrTypes: DtcLbdnPoolsAttrTypes}},
	"priority":                   types.Int64Type,
	"topology":                   types.StringType,
	"ttl":                        types.Int64Type,
	"types":                      internaltypes.UnorderedListOfStringType,
	"use_ttl":                    types.BoolType,
}

var DtcLbdnResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"auth_zones": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "List of linked auth zones.",
	},
	"auto_consolidated_monitors": schema.BoolAttribute{
		Computed:            true,
		Optional:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Flag for enabling auto managing DTC Consolidated Monitors on related DTC Pools.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Comment for the DTC LBDN; maximum 256 characters.",
	},
	"disable": schema.BoolAttribute{
		Computed:            true,
		Optional:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines whether the DTC LBDN is disabled or not. When this is set to False, the fixed address is enabled.",
	},
	"extattrs": schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Extensible attributes associated with the object.",
	},
	"extattrs_all": schema.MapAttribute{
		Computed:            true,
		ElementType:         types.StringType,
		MarkdownDescription: "Extensible attributes associated with the object , including default attributes.",
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
	},
	"health": schema.SingleNestedAttribute{
		Attributes: DtcLbdnHealthResourceSchemaAttributes,
		Computed:   true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The LBDN health information.",
	},
	"lb_method": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("GLOBAL_AVAILABILITY", "RATIO", "ROUND_ROBIN", "SOURCE_IP_HASH", "TOPOLOGY"),
		},
		MarkdownDescription: "The load balancing method. Used to select pool.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The display name of the DTC LBDN, not DNS related.",
	},
	"patterns": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Optional:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "LBDN wildcards for pattern match.",
	},
	"persistence": schema.Int64Attribute{
		Computed:            true,
		Optional:            true,
		Default:             int64default.StaticInt64(0),
		MarkdownDescription: "Maximum time, in seconds, for which client specific LBDN responses will be cached. Zero specifies no caching.",
	},
	"pools": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: DtcLbdnPoolsResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The maximum time, in seconds, for which client specific LBDN responses will be cached. Zero specifies no caching.",
	},
	"priority": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(1),
		MarkdownDescription: "The LBDN pattern match priority for \"overlapping\" DTC LBDN objects. LBDNs are \"overlapping\" if they are simultaneously assigned to a zone and have patterns that can match the same FQDN. The matching LBDN with highest priority (lowest ordinal) will be used.",
	},
	"topology": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The topology rules for TOPOLOGY method.",
	},
	"ttl": schema.Int64Attribute{
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRoot("use_ttl")),
		},
		MarkdownDescription: "Time-to-live value of the record, in seconds.",
	},
	"types": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			customvalidator.StringsInSlice([]string{"A", "AAAA", "CNAME", "NAPTR", "SRV"}),
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The list of resource record types supported by LBDN.",
	},
	"use_ttl": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Flag to indicate whether the TTL value should be used for the LBDN record.",
	},
}

func (m *DtcLbdnModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dtc.DtcLbdn {
	if m == nil {
		return nil
	}
	to := &dtc.DtcLbdn{
		AuthZones:                flex.ExpandFrameworkListString(ctx, m.AuthZones, diags),
		AutoConsolidatedMonitors: flex.ExpandBoolPointer(m.AutoConsolidatedMonitors),
		Comment:                  flex.ExpandStringPointer(m.Comment),
		Disable:                  flex.ExpandBoolPointer(m.Disable),
		ExtAttrs:                 ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		LbMethod:                 flex.ExpandStringPointer(m.LbMethod),
		Name:                     flex.ExpandStringPointer(m.Name),
		Patterns:                 flex.ExpandFrameworkListString(ctx, m.Patterns, diags),
		Persistence:              flex.ExpandInt64Pointer(m.Persistence),
		Pools:                    flex.ExpandFrameworkListNestedBlock(ctx, m.Pools, diags, ExpandDtcLbdnPools),
		Priority:                 flex.ExpandInt64Pointer(m.Priority),
		Topology:                 flex.ExpandStringPointer(m.Topology),
		Ttl:                      flex.ExpandInt64Pointer(m.Ttl),
		Types:                    flex.ExpandFrameworkListString(ctx, m.Types, diags),
		UseTtl:                   flex.ExpandBoolPointer(m.UseTtl),
	}
	return to
}

func FlattenDtcLbdn(ctx context.Context, from *dtc.DtcLbdn, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DtcLbdnAttrTypes)
	}
	m := DtcLbdnModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, DtcLbdnAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DtcLbdnModel) Flatten(ctx context.Context, from *dtc.DtcLbdn, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DtcLbdnModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AuthZones = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.AuthZones, diags)
	m.AutoConsolidatedMonitors = types.BoolPointerValue(from.AutoConsolidatedMonitors)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Health = FlattenDtcLbdnHealth(ctx, from.Health, diags)
	m.LbMethod = flex.FlattenStringPointer(from.LbMethod)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Patterns = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.Patterns, diags)
	m.Persistence = flex.FlattenInt64Pointer(from.Persistence)
	m.Pools = flex.FlattenFrameworkListNestedBlock(ctx, from.Pools, DtcLbdnPoolsAttrTypes, diags, FlattenDtcLbdnPools)
	m.Priority = flex.FlattenInt64Pointer(from.Priority)
	m.Topology = flex.FlattenStringPointer(from.Topology)
	m.Ttl = flex.FlattenInt64Pointer(from.Ttl)
	m.Types = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.Types, diags)
	m.UseTtl = types.BoolPointerValue(from.UseTtl)
}

func (m *DtcLbdnModel) PutExpand(to *dtc.DtcLbdn) *dtc.DtcLbdn {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range DtcLbdnResourceSchemaAttributes {
		attrVal := reflect.ValueOf(attr)
		attrType := attrVal.Type()
		if toType.Kind() == reflect.Struct {
			for i := 0; i < toType.NumField(); i++ {
				fieldValue := toVal.Field(i).Interface()
				tField := toType.Field(i)
				cleanTag := strings.Split(tField.Tag.Get("json"), ",")[0]
				cleanTag = strings.Trim(cleanTag, "_")
				txtFieldValue := utils.ToString(field, fieldValue)
				if field == cleanTag {
					_, ok := attrType.FieldByName("Default")
					if ok {
						defaultVal := attrVal.FieldByName("Default")
						if defaultVal.IsValid() && defaultVal.CanInterface() {
							strDef, ok := defaultVal.Interface().(defaults.String)
							if ok {
								if strDef == stringdefault.StaticString("") {
									continue
								} else if txtFieldValue == "" {
									utils.DeleteBy(to, tField.Name)
								}
							}
							if !ok && txtFieldValue == "" {
								utils.DeleteBy(to, tField.Name)
							}
						}
					} else if txtFieldValue == "" {
						utils.DeleteBy(to, tField.Name)
					}
					_, ok = attrType.FieldByName("Computed")
					if ok {
						computedVal := attrVal.FieldByName("Computed")
						if computedVal.IsValid() && computedVal.CanInterface() {
							boolComp, ok := computedVal.Interface().(bool)
							fmt.Printf("Field: %s, ok: %v, Computed: %v, fieldValue: %v, Value: %s\n", field, ok, boolComp, fieldValue, txtFieldValue)
							if ok {
								if boolComp && txtFieldValue == "" {
									utils.DeleteBy(to, tField.Name)
								}
							} else if txtFieldValue == "" {
								fmt.Printf("Field: %s is marked as computed but is not a bool. Value: %s\n", field, txtFieldValue)
								utils.DeleteBy(to, tField.Name)
							}
						}
					}
					// If the field value is a struct, recursively iterate through its fields
					var deleteEmptyFields func(reflect.Value)
					deleteEmptyFields = func(val reflect.Value) {
						if val.Kind() == reflect.Ptr {
							if val.IsNil() {
								return
							}
							val = val.Elem()
						}
						if val.Kind() != reflect.Struct {
							return
						}
						valType := val.Type()
						for j := 0; j < valType.NumField(); j++ {
							subField := valType.Field(j)
							subFieldValue := val.Field(j)
							subFieldName := strings.Split(subField.Tag.Get("json"), ",")[0]
							subFieldName = strings.Trim(subFieldName, "_")
							txtSubFieldValue := utils.ToString(subFieldName, subFieldValue.Interface())
							if subFieldValue.Kind() == reflect.Struct {
								deleteEmptyFields(subFieldValue)
							}
							if txtSubFieldValue == "" {
								utils.DeleteBy(val.Addr().Interface(), subField.Name)
							}
						}
					}
					if reflect.TypeOf(fieldValue).Kind() == reflect.Struct {
						deleteEmptyFields(reflect.ValueOf(fieldValue))
					} else if reflect.TypeOf(fieldValue).Kind() == reflect.Slice || reflect.TypeOf(fieldValue).Kind() == reflect.Array {
						sliceVal := reflect.ValueOf(fieldValue)
						for i := 0; i < sliceVal.Len(); i++ {
							elem := sliceVal.Index(i)
							if elem.Kind() == reflect.Ptr {
								elem = elem.Elem()
							}
							if elem.Kind() == reflect.Struct {
								deleteEmptyFields(elem)
							}
						}
					}
				}
			}
		}
	}
	return to
}
