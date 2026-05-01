package ipam

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type BulkhostnametemplateModel struct {
	Ref            types.String `tfsdk:"ref"`
	IsGridDefault  types.Bool   `tfsdk:"is_grid_default"`
	PreDefined     types.Bool   `tfsdk:"pre_defined"`
	TemplateFormat types.String `tfsdk:"template_format"`
	TemplateName   types.String `tfsdk:"template_name"`
}

var BulkhostnametemplateAttrTypes = map[string]attr.Type{
	"ref":             types.StringType,
	"is_grid_default": types.BoolType,
	"pre_defined":     types.BoolType,
	"template_format": types.StringType,
	"template_name":   types.StringType,
}

var BulkhostnametemplateResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"is_grid_default": schema.BoolAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "True if this template is Grid default.",
	},
	"pre_defined": schema.BoolAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "True if this is a pre-defined template, False otherwise.",
	},
	"template_format": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(`.*(\$4|#4).*`),
				"Template format must contain at least one of $4 or #4 placeholders",
			),
		},
		MarkdownDescription: "The format of bulk host name template. It should follow certain rules (please use Administration Guide as reference).",
	},
	"template_name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of bulk host name template.",
	},
}

func ExpandBulkhostnametemplate(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.Bulkhostnametemplate {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m BulkhostnametemplateModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *BulkhostnametemplateModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.Bulkhostnametemplate {
	if m == nil {
		return nil
	}
	to := &ipam.Bulkhostnametemplate{
		TemplateFormat: flex.ExpandStringPointer(m.TemplateFormat),
		TemplateName:   flex.ExpandStringPointer(m.TemplateName),
	}
	return to
}

func FlattenBulkhostnametemplate(ctx context.Context, from *ipam.Bulkhostnametemplate, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(BulkhostnametemplateAttrTypes)
	}
	m := BulkhostnametemplateModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, BulkhostnametemplateAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *BulkhostnametemplateModel) Flatten(ctx context.Context, from *ipam.Bulkhostnametemplate, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = BulkhostnametemplateModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.IsGridDefault = types.BoolPointerValue(from.IsGridDefault)
	m.PreDefined = types.BoolPointerValue(from.PreDefined)
	m.TemplateFormat = flex.FlattenStringPointer(from.TemplateFormat)
	m.TemplateName = flex.FlattenStringPointer(from.TemplateName)
}

func (m *BulkhostnametemplateModel) PutExpand(to *ipam.Bulkhostnametemplate) *ipam.Bulkhostnametemplate {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range BulkhostnametemplateResourceSchemaAttributes {
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
