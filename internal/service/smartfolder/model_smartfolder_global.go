package smartfolder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/smartfolder"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type SmartfolderGlobalModel struct {
	Ref        types.String `tfsdk:"ref"`
	Comment    types.String `tfsdk:"comment"`
	GroupBys   types.List   `tfsdk:"group_bys"`
	Name       types.String `tfsdk:"name"`
	QueryItems types.List   `tfsdk:"query_items"`
}

var SmartfolderGlobalAttrTypes = map[string]attr.Type{
	"ref":         types.StringType,
	"comment":     types.StringType,
	"group_bys":   types.ListType{ElemType: types.ObjectType{AttrTypes: SmartfolderGlobalGroupBysAttrTypes}},
	"name":        types.StringType,
	"query_items": types.ListType{ElemType: types.ObjectType{AttrTypes: SmartfolderGlobalQueryItemsAttrTypes}},
}

var SmartfolderGlobalResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The global Smart Folder descriptive comment.",
	},
	"group_bys": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: SmartfolderGlobalGroupBysResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
		MarkdownDescription: "Global Smart Folder grouping rules.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The global Smart Folder name.",
	},
	"query_items": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: SmartfolderGlobalQueryItemsResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		Default: listdefault.StaticValue(types.ListValueMust(
			types.ObjectType{AttrTypes: SmartfolderGlobalQueryItemsAttrTypes},
			[]attr.Value{
				types.ObjectValueMust(SmartfolderGlobalQueryItemsAttrTypes, map[string]attr.Value{
					"name":       types.StringValue("type"),
					"field_type": types.StringValue("NORMAL"),
					"operator":   types.StringValue("EQ"),
					"op_match":   types.BoolValue(true),
					"value_type": types.StringValue("ENUM"),
					"value": types.ObjectValueMust(SmartfolderglobalqueryitemsValueAttrTypes, map[string]attr.Value{
						"value_string":  types.StringValue("Network/Zone/Range/Member"),
						"value_integer": types.Int64Null(),
						"value_date":    types.StringNull(),
						"value_boolean": types.BoolNull(),
					}),
				}),
			},
		)),
		Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
		MarkdownDescription: "The global Smart Folder filter queries.",
	},
}

func (m *SmartfolderGlobalModel) Expand(ctx context.Context, diags *diag.Diagnostics) *smartfolder.SmartfolderGlobal {
	if m == nil {
		return nil
	}
	to := &smartfolder.SmartfolderGlobal{
		Comment:    flex.ExpandStringPointer(m.Comment),
		GroupBys:   flex.ExpandFrameworkListNestedBlock(ctx, m.GroupBys, diags, ExpandSmartfolderGlobalGroupBys),
		Name:       flex.ExpandStringPointer(m.Name),
		QueryItems: flex.ExpandFrameworkListNestedBlock(ctx, m.QueryItems, diags, ExpandSmartfolderGlobalQueryItems),
	}
	return to
}

func FlattenSmartfolderGlobal(ctx context.Context, from *smartfolder.SmartfolderGlobal, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SmartfolderGlobalAttrTypes)
	}
	m := SmartfolderGlobalModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, SmartfolderGlobalAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SmartfolderGlobalModel) Flatten(ctx context.Context, from *smartfolder.SmartfolderGlobal, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SmartfolderGlobalModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.GroupBys = flex.FlattenFrameworkListNestedBlock(ctx, from.GroupBys, SmartfolderGlobalGroupBysAttrTypes, diags, FlattenSmartfolderGlobalGroupBys)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.QueryItems = flex.FlattenFrameworkListNestedBlock(ctx, from.QueryItems, SmartfolderGlobalQueryItemsAttrTypes, diags, FlattenSmartfolderGlobalQueryItems)
}

func (m *SmartfolderGlobalModel) PutExpand(to *smartfolder.SmartfolderGlobal) *smartfolder.SmartfolderGlobal {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range SmartfolderGlobalResourceSchemaAttributes {
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
