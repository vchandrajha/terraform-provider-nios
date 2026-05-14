package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupSamlSettingModel struct {
	AutoCreateUser         types.Bool `tfsdk:"auto_create_user"`
	PersistAutoCreatedUser types.Bool `tfsdk:"persist_auto_created_user"`
}

var AdmingroupSamlSettingAttrTypes = map[string]attr.Type{
	"auto_create_user":          types.BoolType,
	"persist_auto_created_user": types.BoolType,
}

var AdmingroupSamlSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"auto_create_user": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Whether or not auto create user.",
	},
	"persist_auto_created_user": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Whether or not persist auto created user after logout.",
	},
}

func ExpandAdmingroupSamlSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupSamlSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupSamlSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupSamlSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupSamlSetting {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupSamlSetting{
		AutoCreateUser:         flex.ExpandBoolPointer(m.AutoCreateUser),
		PersistAutoCreatedUser: flex.ExpandBoolPointer(m.PersistAutoCreatedUser),
	}
	return to
}

func FlattenAdmingroupSamlSetting(ctx context.Context, from *security.AdmingroupSamlSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupSamlSettingAttrTypes)
	}
	m := AdmingroupSamlSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupSamlSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupSamlSettingModel) Flatten(ctx context.Context, from *security.AdmingroupSamlSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupSamlSettingModel{}
	}
	m.AutoCreateUser = types.BoolPointerValue(from.AutoCreateUser)
	m.PersistAutoCreatedUser = types.BoolPointerValue(from.PersistAutoCreatedUser)
}

func (m *AdmingroupSamlSettingModel) PutExpand(to *security.AdmingroupSamlSetting) *security.AdmingroupSamlSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AdmingroupSamlSettingResourceSchemaAttributes {
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
