package notification

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

	"github.com/infobloxopen/infoblox-nios-go-client/notification"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NotificationruletemplateinstanceParametersModel struct {
	Name         types.String `tfsdk:"name"`
	Value        types.String `tfsdk:"value"`
	DefaultValue types.String `tfsdk:"default_value"`
	Syntax       types.String `tfsdk:"syntax"`
}

var NotificationruletemplateinstanceParametersAttrTypes = map[string]attr.Type{
	"name":          types.StringType,
	"value":         types.StringType,
	"default_value": types.StringType,
	"syntax":        types.StringType,
}

var NotificationruletemplateinstanceParametersResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name of the REST API template parameter.",
	},
	"value": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The value of the REST API template parameter.",
	},
	"default_value": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The default value of the REST API template parameter.",
	},
	"syntax": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("BOOL", "INT", "STR"),
		},
		MarkdownDescription: "The syntax of the REST API template parameter.",
	},
}

func ExpandNotificationruletemplateinstanceParameters(ctx context.Context, o types.Object, diags *diag.Diagnostics) *notification.NotificationruletemplateinstanceParameters {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NotificationruletemplateinstanceParametersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NotificationruletemplateinstanceParametersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *notification.NotificationruletemplateinstanceParameters {
	if m == nil {
		return nil
	}
	to := &notification.NotificationruletemplateinstanceParameters{
		Name:   flex.ExpandStringPointer(m.Name),
		Value:  flex.ExpandStringPointer(m.Value),
		Syntax: flex.ExpandStringPointer(m.Syntax),
	}
	return to
}

func FlattenNotificationruletemplateinstanceParameters(ctx context.Context, from *notification.NotificationruletemplateinstanceParameters, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NotificationruletemplateinstanceParametersAttrTypes)
	}
	m := NotificationruletemplateinstanceParametersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NotificationruletemplateinstanceParametersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NotificationruletemplateinstanceParametersModel) Flatten(ctx context.Context, from *notification.NotificationruletemplateinstanceParameters, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NotificationruletemplateinstanceParametersModel{}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Value = flex.FlattenStringPointer(from.Value)
	m.DefaultValue = flex.FlattenStringPointer(from.DefaultValue)
	m.Syntax = flex.FlattenStringPointer(from.Syntax)
}

func (m *NotificationruletemplateinstanceParametersModel) PutExpand(to *notification.NotificationruletemplateinstanceParameters) *notification.NotificationruletemplateinstanceParameters {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NotificationruletemplateinstanceParametersResourceSchemaAttributes {
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
