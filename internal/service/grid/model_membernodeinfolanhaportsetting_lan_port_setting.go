package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MembernodeinfolanhaportsettingLanPortSettingModel struct {
	AutoPortSettingEnabled types.Bool   `tfsdk:"auto_port_setting_enabled"`
	Speed                  types.String `tfsdk:"speed"`
	Duplex                 types.String `tfsdk:"duplex"`
}

var MembernodeinfolanhaportsettingLanPortSettingAttrTypes = map[string]attr.Type{
	"auto_port_setting_enabled": types.BoolType,
	"speed":                     types.StringType,
	"duplex":                    types.StringType,
}

var MembernodeinfolanhaportsettingLanPortSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"auto_port_setting_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Enable or disable the auto port setting.",
	},
	"speed": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The port speed; if speed is 1000, duplex is FULL.",
	},
	"duplex": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The port duplex; if speed is 1000, duplex must be FULL.",
	},
}

func ExpandMembernodeinfolanhaportsettingLanPortSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MembernodeinfolanhaportsettingLanPortSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MembernodeinfolanhaportsettingLanPortSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MembernodeinfolanhaportsettingLanPortSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MembernodeinfolanhaportsettingLanPortSetting {
	if m == nil {
		return nil
	}
	to := &grid.MembernodeinfolanhaportsettingLanPortSetting{
		AutoPortSettingEnabled: flex.ExpandBoolPointer(m.AutoPortSettingEnabled),
		Speed:                  flex.ExpandStringPointerEmptyAsNil(m.Speed),
		Duplex:                 flex.ExpandStringPointerEmptyAsNil(m.Duplex),
	}
	return to
}

func FlattenMembernodeinfolanhaportsettingLanPortSetting(ctx context.Context, from *grid.MembernodeinfolanhaportsettingLanPortSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MembernodeinfolanhaportsettingLanPortSettingAttrTypes)
	}
	m := MembernodeinfolanhaportsettingLanPortSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MembernodeinfolanhaportsettingLanPortSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MembernodeinfolanhaportsettingLanPortSettingModel) Flatten(ctx context.Context, from *grid.MembernodeinfolanhaportsettingLanPortSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MembernodeinfolanhaportsettingLanPortSettingModel{}
	}
	m.AutoPortSettingEnabled = types.BoolPointerValue(from.AutoPortSettingEnabled)
	m.Speed = flex.FlattenStringPointer(from.Speed)
	m.Duplex = flex.FlattenStringPointer(from.Duplex)
}

func (m *MembernodeinfolanhaportsettingLanPortSettingModel) PutExpand(to *grid.MembernodeinfolanhaportsettingLanPortSetting) *grid.MembernodeinfolanhaportsettingLanPortSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MembernodeinfolanhaportsettingLanPortSettingResourceSchemaAttributes {
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
