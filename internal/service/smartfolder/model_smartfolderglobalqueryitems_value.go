package smartfolder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/smartfolder"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type SmartfolderglobalqueryitemsValueModel struct {
	ValueInteger types.Int64  `tfsdk:"value_integer"`
	ValueString  types.String `tfsdk:"value_string"`
	ValueDate    types.String `tfsdk:"value_date"`
	ValueBoolean types.Bool   `tfsdk:"value_boolean"`
}

var SmartfolderglobalqueryitemsValueAttrTypes = map[string]attr.Type{
	"value_integer": types.Int64Type,
	"value_string":  types.StringType,
	"value_date":    types.StringType,
	"value_boolean": types.BoolType,
}

var SmartfolderglobalqueryitemsValueResourceSchemaAttributes = map[string]schema.Attribute{
	"value_integer": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "The integer value of the Smart Folder query.",
	},
	"value_string": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The string value of the Smart Folder query.",
	},
	"value_date": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTimeFormat(),
		},
		MarkdownDescription: "The timestamp value of the Smart Folder query.",
	},
	"value_boolean": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "The boolean value of the Smart Folder query.",
	},
}

func ExpandSmartfolderglobalqueryitemsValue(ctx context.Context, o types.Object, diags *diag.Diagnostics) *smartfolder.SmartfolderglobalqueryitemsValue {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m SmartfolderglobalqueryitemsValueModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *SmartfolderglobalqueryitemsValueModel) Expand(ctx context.Context, diags *diag.Diagnostics) *smartfolder.SmartfolderglobalqueryitemsValue {
	if m == nil {
		return nil
	}
	to := &smartfolder.SmartfolderglobalqueryitemsValue{
		ValueInteger: flex.ExpandInt64Pointer(m.ValueInteger),
		ValueString:  flex.ExpandStringPointer(m.ValueString),
		ValueBoolean: flex.ExpandBoolPointer(m.ValueBoolean),
	}
	to.ValueDate = flex.ExpandTimeToUnix(m.ValueDate, diags)
	return to
}

func FlattenSmartfolderglobalqueryitemsValue(ctx context.Context, from *smartfolder.SmartfolderglobalqueryitemsValue, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SmartfolderglobalqueryitemsValueAttrTypes)
	}
	m := SmartfolderglobalqueryitemsValueModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, SmartfolderglobalqueryitemsValueAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SmartfolderglobalqueryitemsValueModel) Flatten(ctx context.Context, from *smartfolder.SmartfolderglobalqueryitemsValue, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SmartfolderglobalqueryitemsValueModel{}
	}
	m.ValueInteger = flex.FlattenInt64Pointer(from.ValueInteger)
	if from.ValueString != nil {
		m.ValueString = flex.FlattenStringPointer(from.ValueString)
	} else {
		m.ValueString = types.StringNull()
	}
	if from.ValueDate != nil {
		m.ValueDate = flex.FlattenUnixTime(from.ValueDate, diags)
	} else {
		m.ValueDate = types.StringNull()
	}
	m.ValueBoolean = types.BoolPointerValue(from.ValueBoolean)
}

func (m *SmartfolderglobalqueryitemsValueModel) PutExpand(to *smartfolder.SmartfolderglobalqueryitemsValue) *smartfolder.SmartfolderglobalqueryitemsValue {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range SmartfolderglobalqueryitemsValueResourceSchemaAttributes {
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
