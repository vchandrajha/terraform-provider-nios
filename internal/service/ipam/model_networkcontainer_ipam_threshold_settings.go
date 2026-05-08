package ipam

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NetworkcontainerIpamThresholdSettingsModel struct {
	TriggerValue types.Int64 `tfsdk:"trigger_value"`
	ResetValue   types.Int64 `tfsdk:"reset_value"`
}

var NetworkcontainerIpamThresholdSettingsAttrTypes = map[string]attr.Type{
	"trigger_value": types.Int64Type,
	"reset_value":   types.Int64Type,
}

var NetworkcontainerIpamThresholdSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"trigger_value": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Indicates the percentage point which triggers the email/SNMP trap sending.",
		Default:             int64default.StaticInt64(95),
	},
	"reset_value": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(85),
		MarkdownDescription: "Indicates the percentage point which resets the email/SNMP trap sending.",
	},
}

func ExpandNetworkcontainerIpamThresholdSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworkcontainerIpamThresholdSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworkcontainerIpamThresholdSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworkcontainerIpamThresholdSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworkcontainerIpamThresholdSettings {
	if m == nil {
		return nil
	}
	to := &ipam.NetworkcontainerIpamThresholdSettings{
		TriggerValue: flex.ExpandInt64Pointer(m.TriggerValue),
		ResetValue:   flex.ExpandInt64Pointer(m.ResetValue),
	}
	return to
}

func FlattenNetworkcontainerIpamThresholdSettings(ctx context.Context, from *ipam.NetworkcontainerIpamThresholdSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworkcontainerIpamThresholdSettingsAttrTypes)
	}
	m := NetworkcontainerIpamThresholdSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworkcontainerIpamThresholdSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworkcontainerIpamThresholdSettingsModel) Flatten(ctx context.Context, from *ipam.NetworkcontainerIpamThresholdSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworkcontainerIpamThresholdSettingsModel{}
	}
	m.TriggerValue = flex.FlattenInt64Pointer(from.TriggerValue)
	m.ResetValue = flex.FlattenInt64Pointer(from.ResetValue)
}

func (m *NetworkcontainerIpamThresholdSettingsModel) PutExpand(to *ipam.NetworkcontainerIpamThresholdSettings) *ipam.NetworkcontainerIpamThresholdSettings {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NetworkcontainerIpamThresholdSettingsResourceSchemaAttributes {
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
