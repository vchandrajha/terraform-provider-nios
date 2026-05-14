package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type ViewDnssecTrustedKeysModel struct {
	Fqdn               types.String `tfsdk:"fqdn"`
	Algorithm          types.String `tfsdk:"algorithm"`
	Key                types.String `tfsdk:"key"`
	SecureEntryPoint   types.Bool   `tfsdk:"secure_entry_point"`
	DnssecMustBeSecure types.Bool   `tfsdk:"dnssec_must_be_secure"`
}

var ViewDnssecTrustedKeysAttrTypes = map[string]attr.Type{
	"fqdn":                  types.StringType,
	"algorithm":             types.StringType,
	"key":                   types.StringType,
	"secure_entry_point":    types.BoolType,
	"dnssec_must_be_secure": types.BoolType,
}

var ViewDnssecTrustedKeysResourceSchemaAttributes = map[string]schema.Attribute{
	"fqdn": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The FQDN of the domain for which the member validates responses to recursive queries.",
	},
	"algorithm": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The DNSSEC algorithm used to generate the key.",
	},
	"key": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The DNSSEC key.",
	},
	"secure_entry_point": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The secure entry point flag, if set it means this is a KSK configuration.",
	},
	"dnssec_must_be_secure": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Responses must be DNSSEC secure for this hierarchy/domain.",
	},
}

func ExpandViewDnssecTrustedKeys(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ViewDnssecTrustedKeys {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ViewDnssecTrustedKeysModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ViewDnssecTrustedKeysModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ViewDnssecTrustedKeys {
	if m == nil {
		return nil
	}
	to := &dns.ViewDnssecTrustedKeys{
		Fqdn:               flex.ExpandStringPointer(m.Fqdn),
		Algorithm:          flex.ExpandStringPointer(m.Algorithm),
		Key:                flex.ExpandStringPointer(m.Key),
		SecureEntryPoint:   flex.ExpandBoolPointer(m.SecureEntryPoint),
		DnssecMustBeSecure: flex.ExpandBoolPointer(m.DnssecMustBeSecure),
	}
	return to
}

func FlattenViewDnssecTrustedKeys(ctx context.Context, from *dns.ViewDnssecTrustedKeys, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ViewDnssecTrustedKeysAttrTypes)
	}
	m := ViewDnssecTrustedKeysModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ViewDnssecTrustedKeysAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ViewDnssecTrustedKeysModel) Flatten(ctx context.Context, from *dns.ViewDnssecTrustedKeys, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ViewDnssecTrustedKeysModel{}
	}
	m.Fqdn = flex.FlattenStringPointer(from.Fqdn)
	m.Algorithm = flex.FlattenStringPointer(from.Algorithm)
	m.Key = flex.FlattenStringPointer(from.Key)
	m.SecureEntryPoint = types.BoolPointerValue(from.SecureEntryPoint)
	m.DnssecMustBeSecure = types.BoolPointerValue(from.DnssecMustBeSecure)
}

func (m *ViewDnssecTrustedKeysModel) PutExpand(to *dns.ViewDnssecTrustedKeys) *dns.ViewDnssecTrustedKeys {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ViewDnssecTrustedKeysResourceSchemaAttributes {
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
							if ok {
								if boolComp && txtFieldValue == "" {
									utils.DeleteBy(to, tField.Name)
								}
							} else if txtFieldValue == "" {
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
