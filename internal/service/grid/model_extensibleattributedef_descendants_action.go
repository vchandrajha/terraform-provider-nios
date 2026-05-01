package grid

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

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ExtensibleattributedefDescendantsActionModel struct {
	OptionWithEa    types.String `tfsdk:"option_with_ea"`
	OptionWithoutEa types.String `tfsdk:"option_without_ea"`
	OptionDeleteEa  types.String `tfsdk:"option_delete_ea"`
}

var ExtensibleattributedefDescendantsActionAttrTypes = map[string]attr.Type{
	"option_with_ea":    types.StringType,
	"option_without_ea": types.StringType,
	"option_delete_ea":  types.StringType,
}

var ExtensibleattributedefDescendantsActionResourceSchemaAttributes = map[string]schema.Attribute{
	"option_with_ea": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("CONVERT", "INHERIT", "RETAIN"),
		},
		MarkdownDescription: "This option describes which action must be taken if the extensible attribute exists for both the parent and descendant objects: * INHERIT: inherit the extensible attribute from the parent object. * RETAIN: retain the value of an extensible attribute that was set for the child object. * CONVERT: the value of the extensible attribute must be copied from the parent object.",
	},
	"option_without_ea": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("INHERIT", "NOT_INHERIT"),
		},
		MarkdownDescription: "This option describes which action must be taken if the extensible attribute exists for the parent, but is absent from the descendant object: * INHERIT: inherit the extensible attribute from the parent object. * NOT_INHERIT: do nothing.",
	},
	"option_delete_ea": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("REMOVE", "RETAIN"),
		},
		MarkdownDescription: "This option describes which action must be taken if the extensible attribute exists for the descendant, but is absent for the parent object: * RETAIN: retain the extensible attribute value for the descendant object. * REMOVE: remove this extensible attribute from the descendant object.",
	},
}

func ExpandExtensibleattributedefDescendantsAction(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.ExtensibleattributedefDescendantsAction {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ExtensibleattributedefDescendantsActionModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ExtensibleattributedefDescendantsActionModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.ExtensibleattributedefDescendantsAction {
	if m == nil {
		return nil
	}
	to := &grid.ExtensibleattributedefDescendantsAction{
		OptionWithEa:    flex.ExpandStringPointer(m.OptionWithEa),
		OptionWithoutEa: flex.ExpandStringPointer(m.OptionWithoutEa),
		OptionDeleteEa:  flex.ExpandStringPointer(m.OptionDeleteEa),
	}
	return to
}

func FlattenExtensibleattributedefDescendantsAction(ctx context.Context, from *grid.ExtensibleattributedefDescendantsAction, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ExtensibleattributedefDescendantsActionAttrTypes)
	}
	m := ExtensibleattributedefDescendantsActionModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ExtensibleattributedefDescendantsActionAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ExtensibleattributedefDescendantsActionModel) Flatten(ctx context.Context, from *grid.ExtensibleattributedefDescendantsAction, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ExtensibleattributedefDescendantsActionModel{}
	}
	m.OptionWithEa = flex.FlattenStringPointer(from.OptionWithEa)
	m.OptionWithoutEa = flex.FlattenStringPointer(from.OptionWithoutEa)
	m.OptionDeleteEa = flex.FlattenStringPointer(from.OptionDeleteEa)
}

func (m *ExtensibleattributedefDescendantsActionModel) PutExpand(to *grid.ExtensibleattributedefDescendantsAction) *grid.ExtensibleattributedefDescendantsAction {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ExtensibleattributedefDescendantsActionResourceSchemaAttributes {
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
