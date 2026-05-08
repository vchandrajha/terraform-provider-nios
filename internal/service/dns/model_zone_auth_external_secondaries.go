package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type ZoneAuthExternalSecondariesModel struct {
	Address                      iptypes.IPAddress `tfsdk:"address"`
	Name                         types.String      `tfsdk:"name"`
	SharedWithMsParentDelegation types.Bool        `tfsdk:"shared_with_ms_parent_delegation"`
	Stealth                      types.Bool        `tfsdk:"stealth"`
	TsigKey                      types.String      `tfsdk:"tsig_key"`
	TsigKeyAlg                   types.String      `tfsdk:"tsig_key_alg"`
	TsigKeyName                  types.String      `tfsdk:"tsig_key_name"`
	UseTsigKeyName               types.Bool        `tfsdk:"use_tsig_key_name"`
}

var ZoneAuthExternalSecondariesAttrTypes = map[string]attr.Type{
	"address":                          iptypes.IPAddressType{},
	"name":                             types.StringType,
	"shared_with_ms_parent_delegation": types.BoolType,
	"stealth":                          types.BoolType,
	"tsig_key":                         types.StringType,
	"tsig_key_alg":                     types.StringType,
	"tsig_key_name":                    types.StringType,
	"use_tsig_key_name":                types.BoolType,
}

var ZoneAuthExternalSecondariesResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Required:            true,
		MarkdownDescription: "The IPv4 Address or IPv6 Address of the server.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "A resolvable domain name for the external DNS server.",
	},
	"shared_with_ms_parent_delegation": schema.BoolAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "This flag represents whether the name server is shared with the parent Microsoft primary zone's delegation server.",
	},
	"stealth": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Set this flag to hide the NS record for the primary name server from DNS queries.",
	},
	"tsig_key": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},

		MarkdownDescription: "A generated TSIG key.",
	},
	"tsig_key_alg": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("HMAC-MD5"),
		Validators: []validator.String{
			stringvalidator.OneOf("HMAC-MD5", "HMAC-SHA256"),
		},
		MarkdownDescription: "The TSIG key algorithm.",
	},
	"tsig_key_name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The TSIG key name.",
	},
	"use_tsig_key_name": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Use flag for: tsig_key_name",
	},
}

func ExpandZoneAuthExternalSecondaries(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthExternalSecondaries {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthExternalSecondariesModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthExternalSecondariesModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthExternalSecondaries {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthExternalSecondaries{
		Address:        flex.ExpandIPAddress(m.Address),
		Name:           flex.ExpandStringPointer(m.Name),
		Stealth:        flex.ExpandBoolPointer(m.Stealth),
		TsigKey:        flex.ExpandStringPointer(m.TsigKey),
		TsigKeyAlg:     flex.ExpandStringPointer(m.TsigKeyAlg),
		TsigKeyName:    flex.ExpandStringPointer(m.TsigKeyName),
		UseTsigKeyName: flex.ExpandBoolPointer(m.UseTsigKeyName),
	}
	return to
}

func FlattenZoneAuthExternalSecondaries(ctx context.Context, from *dns.ZoneAuthExternalSecondaries, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthExternalSecondariesAttrTypes)
	}
	m := ZoneAuthExternalSecondariesModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthExternalSecondariesAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthExternalSecondariesModel) Flatten(ctx context.Context, from *dns.ZoneAuthExternalSecondaries, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthExternalSecondariesModel{}
	}
	m.Address = flex.FlattenIPAddress(from.Address)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.SharedWithMsParentDelegation = types.BoolPointerValue(from.SharedWithMsParentDelegation)
	m.Stealth = types.BoolPointerValue(from.Stealth)
	m.TsigKey = flex.FlattenStringPointer(from.TsigKey)
	m.TsigKeyAlg = flex.FlattenStringPointer(from.TsigKeyAlg)
	m.TsigKeyName = flex.FlattenStringPointer(from.TsigKeyName)
	m.UseTsigKeyName = types.BoolPointerValue(from.UseTsigKeyName)
}

func (m *ZoneAuthExternalSecondariesModel) PutExpand(to *dns.ZoneAuthExternalSecondaries) *dns.ZoneAuthExternalSecondaries {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthExternalSecondariesResourceSchemaAttributes {
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
