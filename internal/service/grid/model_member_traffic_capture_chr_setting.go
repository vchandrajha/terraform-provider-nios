package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberTrafficCaptureChrSettingModel struct {
	ChrTriggerEnable       types.Bool  `tfsdk:"chr_trigger_enable"`
	ChrThreshold           types.Int64 `tfsdk:"chr_threshold"`
	ChrReset               types.Int64 `tfsdk:"chr_reset"`
	ChrMinCacheUtilization types.Int64 `tfsdk:"chr_min_cache_utilization"`
}

var MemberTrafficCaptureChrSettingAttrTypes = map[string]attr.Type{
	"chr_trigger_enable":        types.BoolType,
	"chr_threshold":             types.Int64Type,
	"chr_reset":                 types.Int64Type,
	"chr_min_cache_utilization": types.Int64Type,
}

var MemberTrafficCaptureChrSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"chr_trigger_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable triggering automated traffic capture based on cache hit ratio thresholds.",
	},
	"chr_threshold": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "DNS Cache hit ratio threshold(%) below which traffic capture will be triggered.",
	},
	"chr_reset": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "DNS Cache hit ratio threshold(%) above which traffic capture will be triggered.",
	},
	"chr_min_cache_utilization": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Minimum DNS cache utilization threshold(%) for triggering traffic capture based on DNS cache hit ratio.",
	},
}

func ExpandMemberTrafficCaptureChrSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberTrafficCaptureChrSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberTrafficCaptureChrSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberTrafficCaptureChrSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberTrafficCaptureChrSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberTrafficCaptureChrSetting{
		ChrTriggerEnable:       flex.ExpandBoolPointer(m.ChrTriggerEnable),
		ChrThreshold:           flex.ExpandInt64Pointer(m.ChrThreshold),
		ChrReset:               flex.ExpandInt64Pointer(m.ChrReset),
		ChrMinCacheUtilization: flex.ExpandInt64Pointer(m.ChrMinCacheUtilization),
	}
	return to
}

func FlattenMemberTrafficCaptureChrSetting(ctx context.Context, from *grid.MemberTrafficCaptureChrSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberTrafficCaptureChrSettingAttrTypes)
	}
	m := MemberTrafficCaptureChrSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberTrafficCaptureChrSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberTrafficCaptureChrSettingModel) Flatten(ctx context.Context, from *grid.MemberTrafficCaptureChrSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberTrafficCaptureChrSettingModel{}
	}
	m.ChrTriggerEnable = types.BoolPointerValue(from.ChrTriggerEnable)
	m.ChrThreshold = flex.FlattenInt64Pointer(from.ChrThreshold)
	m.ChrReset = flex.FlattenInt64Pointer(from.ChrReset)
	m.ChrMinCacheUtilization = flex.FlattenInt64Pointer(from.ChrMinCacheUtilization)
}

func (m *MemberTrafficCaptureChrSettingModel) PutExpand(to *grid.MemberTrafficCaptureChrSetting) *grid.MemberTrafficCaptureChrSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberTrafficCaptureChrSettingResourceSchemaAttributes {
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
