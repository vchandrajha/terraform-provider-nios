package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/infoblox-nios-go-client/ipam"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type FuncCallModel struct {
	AttributeName    types.String `tfsdk:"attribute_name"`
	ObjectFunction   types.String `tfsdk:"object_function"`
	Parameters       types.Map    `tfsdk:"parameters"`
	ResultField      types.String `tfsdk:"result_field"`
	Object           types.String `tfsdk:"object"`
	ObjectParameters types.Map    `tfsdk:"object_parameters"`
}

var FuncCallAttrTypes = map[string]attr.Type{
	"attribute_name":    types.StringType,
	"object_function":   types.StringType,
	"parameters":        types.MapType{ElemType: types.StringType},
	"result_field":      types.StringType,
	"object":            types.StringType,
	"object_parameters": types.MapType{ElemType: types.StringType},
}

var FuncCallResourceSchemaAttributes = map[string]schema.Attribute{
	"attribute_name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The attribute to be called.",
	},
	"object_function": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The function to be called.",
	},
	"parameters": schema.MapAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		MarkdownDescription: "The parameters for the function.",
	},
	"result_field": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The result field of the function.",
	},
	"object": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The object to be called.",
	},
	"object_parameters": schema.MapAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		MarkdownDescription: "The parameters for the object.",
	},
}

func ExpandFuncCall(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.FuncCall {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m FuncCallModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *FuncCallModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.FuncCall {
	if m == nil {
		return nil
	}
	to := &ipam.FuncCall{
		AttributeName:    flex.ExpandString(m.AttributeName),
		ObjectFunction:   flex.ExpandStringPointer(m.ObjectFunction),
		Parameters:       flex.ExpandParsedFrameworkMapString(ctx, m.Parameters, diags),
		ResultField:      flex.ExpandStringPointer(m.ResultField),
		Object:           flex.ExpandStringPointer(m.Object),
		ObjectParameters: flex.ExpandFrameworkMapString(ctx, m.ObjectParameters, diags),
	}
	return to
}

func FlattenFuncCall(ctx context.Context, from *ipam.FuncCall, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(FuncCallAttrTypes)
	}
	m := FuncCallModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, FuncCallAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *FuncCallModel) Flatten(ctx context.Context, from *ipam.FuncCall, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = FuncCallModel{}
	}
	m.AttributeName = flex.FlattenString(from.AttributeName)
	m.ObjectFunction = flex.FlattenStringPointer(from.ObjectFunction)
	m.Parameters = flex.FlattenFrameworkMapString(ctx, from.Parameters, diags)
	m.ResultField = flex.FlattenStringPointer(from.ResultField)
	m.Object = flex.FlattenStringPointer(from.Object)
	m.ObjectParameters = flex.FlattenFrameworkMapString(ctx, from.ObjectParameters, diags)
}

func (m *FuncCallModel) PutExpand(to *ipam.FuncCall) *ipam.FuncCall {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range FuncCallResourceSchemaAttributes {
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
