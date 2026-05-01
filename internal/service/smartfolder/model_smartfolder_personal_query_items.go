package smartfolder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/smartfolder"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type SmartfolderPersonalQueryItemsModel struct {
	Name      types.String `tfsdk:"name"`
	FieldType types.String `tfsdk:"field_type"`
	Operator  types.String `tfsdk:"operator"`
	OpMatch   types.Bool   `tfsdk:"op_match"`
	ValueType types.String `tfsdk:"value_type"`
	Value     types.Object `tfsdk:"value"`
}

var SmartfolderPersonalQueryItemsAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"field_type": types.StringType,
	"operator":   types.StringType,
	"op_match":   types.BoolType,
	"value_type": types.StringType,
	"value":      types.ObjectType{AttrTypes: SmartfolderpersonalqueryitemsValueAttrTypes},
}

var SmartfolderPersonalQueryItemsResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("type"),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The Smart Folder query name.",
	},
	"field_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("NORMAL"),
		Validators: []validator.String{
			stringvalidator.OneOf("EXTATTR", "NORMAL"),
		},
		MarkdownDescription: "The Smart Folder query field type.",
	},
	"operator": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("EQ"),
		Validators: []validator.String{
			stringvalidator.OneOf(
				"BEGINS_WITH",
				"CONTAINS",
				"DROPS_BY",
				"ENDS_WITH",
				"EQ",
				"GEQ",
				"GT",
				"HAS_VALUE",
				"INHERITANCE_STATE_EQUALS",
				"IP_ADDR_WITHIN",
				"LEQ",
				"LT",
				"MATCH_EXPR",
				"RELATIVE_DATE",
				"RISES_BY",
				"SUFFIX_MATCH",
			),
		},
		MarkdownDescription: "The Smart Folder operator used in query.",
	},
	"op_match": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "Determines whether the query operator should match.",
	},
	"value_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf(
				"BOOLEAN",
				"DATE",
				"EMAIL",
				"ENUM",
				"INTEGER",
				"OBJTYPE",
				"STRING",
				"URL",
			),
		},
		Default:             stringdefault.StaticString("ENUM"),
		MarkdownDescription: "The Smart Folder query value type.",
	},
	"value": schema.SingleNestedAttribute{
		Attributes: SmartfolderpersonalqueryitemsValueResourceSchemaAttributes,
		Optional:   true,
		Computed:   true,
		Default: objectdefault.StaticValue(types.ObjectValueMust(
			SmartfolderpersonalqueryitemsValueAttrTypes,
			map[string]attr.Value{
				"value_string":  types.StringValue("Network/Zone/Range/Member"),
				"value_integer": types.Int64Null(),
				"value_date":    types.StringNull(),
				"value_boolean": types.BoolNull(),
			},
		)),
		MarkdownDescription: "The Smart Folder query value.",
	},
}

func ExpandSmartfolderPersonalQueryItems(ctx context.Context, o types.Object, diags *diag.Diagnostics) *smartfolder.SmartfolderPersonalQueryItems {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m SmartfolderPersonalQueryItemsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *SmartfolderPersonalQueryItemsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *smartfolder.SmartfolderPersonalQueryItems {
	if m == nil {
		return nil
	}
	to := &smartfolder.SmartfolderPersonalQueryItems{
		Name:      flex.ExpandStringPointer(m.Name),
		FieldType: flex.ExpandStringPointer(m.FieldType),
		Operator:  flex.ExpandStringPointer(m.Operator),
		OpMatch:   flex.ExpandBoolPointer(m.OpMatch),
		ValueType: flex.ExpandStringPointer(m.ValueType),
		Value:     ExpandSmartfolderpersonalqueryitemsValue(ctx, m.Value, diags),
	}
	return to
}

func FlattenSmartfolderPersonalQueryItems(ctx context.Context, from *smartfolder.SmartfolderPersonalQueryItems, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SmartfolderPersonalQueryItemsAttrTypes)
	}
	m := SmartfolderPersonalQueryItemsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, SmartfolderPersonalQueryItemsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SmartfolderPersonalQueryItemsModel) Flatten(ctx context.Context, from *smartfolder.SmartfolderPersonalQueryItems, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SmartfolderPersonalQueryItemsModel{}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.FieldType = flex.FlattenStringPointer(from.FieldType)
	m.Operator = flex.FlattenStringPointer(from.Operator)
	m.OpMatch = types.BoolPointerValue(from.OpMatch)
	m.ValueType = flex.FlattenStringPointer(from.ValueType)
	m.Value = FlattenSmartfolderpersonalqueryitemsValue(ctx, from.Value, diags)
}

func (m *SmartfolderPersonalQueryItemsModel) PutExpand(to *smartfolder.SmartfolderPersonalQueryItems) *smartfolder.SmartfolderPersonalQueryItems {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range SmartfolderPersonalQueryItemsResourceSchemaAttributes {
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
