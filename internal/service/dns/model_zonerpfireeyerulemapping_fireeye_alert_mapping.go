package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZonerpfireeyerulemappingFireeyeAlertMappingModel struct {
	AlertType types.String `tfsdk:"alert_type"`
	RpzRule   types.String `tfsdk:"rpz_rule"`
	Lifetime  types.Int64  `tfsdk:"lifetime"`
}

var ZonerpfireeyerulemappingFireeyeAlertMappingAttrTypes = map[string]attr.Type{
	"alert_type": types.StringType,
	"rpz_rule":   types.StringType,
	"lifetime":   types.Int64Type,
}

var ZonerpfireeyerulemappingFireeyeAlertMappingResourceSchemaAttributes = map[string]schema.Attribute{
	"alert_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("DOMAIN_MATCH", "INFECTION_MATCH", "MALWARE_CALLBACK", "MALWARE_OBJECT", "WEB_INFECTION"),
		},
		MarkdownDescription: "The type of Fireeye Alert.",
	},
	"rpz_rule": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("NODATA", "NONE", "NXDOMAIN", "PASSTHRU", "SUBSTITUTE"),
		},
		MarkdownDescription: "The RPZ rule for the alert.",
	},
	"lifetime": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The expiration Lifetime of alert type. The 32-bit unsigned integer represents the amount of seconds this alert type will live for. 0 means the alert will never expire.",
	},
}

func ExpandZonerpfireeyerulemappingFireeyeAlertMapping(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZonerpfireeyerulemappingFireeyeAlertMapping {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZonerpfireeyerulemappingFireeyeAlertMappingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZonerpfireeyerulemappingFireeyeAlertMappingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZonerpfireeyerulemappingFireeyeAlertMapping {
	if m == nil {
		return nil
	}
	to := &dns.ZonerpfireeyerulemappingFireeyeAlertMapping{
		AlertType: flex.ExpandStringPointer(m.AlertType),
		RpzRule:   flex.ExpandStringPointer(m.RpzRule),
		Lifetime:  flex.ExpandInt64Pointer(m.Lifetime),
	}
	return to
}

func FlattenZonerpfireeyerulemappingFireeyeAlertMapping(ctx context.Context, from *dns.ZonerpfireeyerulemappingFireeyeAlertMapping, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZonerpfireeyerulemappingFireeyeAlertMappingAttrTypes)
	}
	m := ZonerpfireeyerulemappingFireeyeAlertMappingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZonerpfireeyerulemappingFireeyeAlertMappingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZonerpfireeyerulemappingFireeyeAlertMappingModel) Flatten(ctx context.Context, from *dns.ZonerpfireeyerulemappingFireeyeAlertMapping, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZonerpfireeyerulemappingFireeyeAlertMappingModel{}
	}
	m.AlertType = flex.FlattenStringPointer(from.AlertType)
	m.RpzRule = flex.FlattenStringPointer(from.RpzRule)
	m.Lifetime = flex.FlattenInt64Pointer(from.Lifetime)
}

func (m *ZonerpfireeyerulemappingFireeyeAlertMappingModel) PutExpand(to *dns.ZonerpfireeyerulemappingFireeyeAlertMapping) *dns.ZonerpfireeyerulemappingFireeyeAlertMapping {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZonerpfireeyerulemappingFireeyeAlertMappingResourceSchemaAttributes {
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
