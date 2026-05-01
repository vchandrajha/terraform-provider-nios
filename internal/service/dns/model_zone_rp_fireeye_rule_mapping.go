package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneRpFireeyeRuleMappingModel struct {
	AptOverride           types.String `tfsdk:"apt_override"`
	FireeyeAlertMapping   types.List   `tfsdk:"fireeye_alert_mapping"`
	SubstitutedDomainName types.String `tfsdk:"substituted_domain_name"`
}

var ZoneRpFireeyeRuleMappingAttrTypes = map[string]attr.Type{
	"apt_override":            types.StringType,
	"fireeye_alert_mapping":   types.ListType{ElemType: types.ObjectType{AttrTypes: ZonerpfireeyerulemappingFireeyeAlertMappingAttrTypes}},
	"substituted_domain_name": types.StringType,
}

var ZoneRpFireeyeRuleMappingResourceSchemaAttributes = map[string]schema.Attribute{
	"apt_override": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("NODATA", "NOOVERRIDE", "NXDOMAIN", "PASSTHRU", "SUBSTITUTE"),
		},
		MarkdownDescription: "The override setting for APT alerts.",
	},
	"fireeye_alert_mapping": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZonerpfireeyerulemappingFireeyeAlertMappingResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The FireEye alert mapping.",
	},
	"substituted_domain_name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The domain name to be substituted, this is applicable only when apt_override is set to \"SUBSTITUTE\".",
	},
}

func ExpandZoneRpFireeyeRuleMapping(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneRpFireeyeRuleMapping {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneRpFireeyeRuleMappingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneRpFireeyeRuleMappingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneRpFireeyeRuleMapping {
	if m == nil {
		return nil
	}
	to := &dns.ZoneRpFireeyeRuleMapping{
		AptOverride:           flex.ExpandStringPointer(m.AptOverride),
		FireeyeAlertMapping:   flex.ExpandFrameworkListNestedBlock(ctx, m.FireeyeAlertMapping, diags, ExpandZonerpfireeyerulemappingFireeyeAlertMapping),
		SubstitutedDomainName: flex.ExpandStringPointer(m.SubstitutedDomainName),
	}
	return to
}

func FlattenZoneRpFireeyeRuleMapping(ctx context.Context, from *dns.ZoneRpFireeyeRuleMapping, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneRpFireeyeRuleMappingAttrTypes)
	}
	m := ZoneRpFireeyeRuleMappingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneRpFireeyeRuleMappingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneRpFireeyeRuleMappingModel) Flatten(ctx context.Context, from *dns.ZoneRpFireeyeRuleMapping, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneRpFireeyeRuleMappingModel{}
	}
	m.AptOverride = flex.FlattenStringPointer(from.AptOverride)
	m.FireeyeAlertMapping = flex.FlattenFrameworkListNestedBlock(ctx, from.FireeyeAlertMapping, ZonerpfireeyerulemappingFireeyeAlertMappingAttrTypes, diags, FlattenZonerpfireeyerulemappingFireeyeAlertMapping)
	m.SubstitutedDomainName = flex.FlattenStringPointer(from.SubstitutedDomainName)
}

func (m *ZoneRpFireeyeRuleMappingModel) PutExpand(to *dns.ZoneRpFireeyeRuleMapping) *dns.ZoneRpFireeyeRuleMapping {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneRpFireeyeRuleMappingResourceSchemaAttributes {
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
							fmt.Printf("Field: %s, Computed: %v, fieldValue: %v, Value: %s\n", field, boolComp, fieldValue, txtFieldValue)
							if ok {
								if !boolComp {
									continue
								} else if txtFieldValue == "" {
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
