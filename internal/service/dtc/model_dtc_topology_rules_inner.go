package dtc

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dtc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type DtcTopologyRulesInnerModel struct {
	DestType        types.String `tfsdk:"dest_type"`
	DestinationLink types.String `tfsdk:"destination_link"`
	ReturnType      types.String `tfsdk:"return_type"`
	Topology        types.String `tfsdk:"topology"`
	Valid           types.Bool   `tfsdk:"valid"`
	Sources         types.List   `tfsdk:"sources"`
}

var DtcTopologyRulesInnerAttrTypes = map[string]attr.Type{
	"dest_type":        types.StringType,
	"destination_link": types.StringType,
	"return_type":      types.StringType,
	"topology":         types.StringType,
	"valid":            types.BoolType,
	"sources":          types.ListType{ElemType: types.ObjectType{AttrTypes: DtcTopologyRulesInnerOneOf1SourcesInnerAttrTypes}},
}

var DtcTopologyRulesInnerResourceSchemaAttributes = map[string]schema.Attribute{
	"dest_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("POOL", "SERVER"),
		},
		MarkdownDescription: "The type of the destination for this rule.",
	},
	"destination_link": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The reference to the destination object.",
	},
	"return_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("NOERR", "NXDOMAIN", "REGULAR"),
		},
		MarkdownDescription: "The type of the return value for this source.",
	},
	"topology": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The topology for this rule.",
	},
	"valid": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Indicates whether the rule is valid.",
	},
	"sources": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: DtcTopologyRulesInnerOneOf1SourcesInnerResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Conditions for matching sources.",
	},
}

func ExpandDtcTopologyRulesInner(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dtc.DtcTopologyRulesInner {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m DtcTopologyRulesInnerModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *DtcTopologyRulesInnerModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dtc.DtcTopologyRulesInner {
	if m == nil {
		return nil
	}
	to := &dtc.DtcTopologyRulesInner{
		DtcTopologyRulesInnerOneOf1: &dtc.DtcTopologyRulesInnerOneOf1{
			DestType:        flex.ExpandStringPointer(m.DestType),
			DestinationLink: flex.ExpandStringPointer(m.DestinationLink),
			ReturnType:      flex.ExpandStringPointer(m.ReturnType),
			Sources:         flex.ExpandFrameworkListNestedBlock(ctx, m.Sources, diags, ExpandDtcTopologyRulesInnerOneOf1SourcesInner),
		},
	}
	return to
}

func FlattenDtcTopologyRulesInner(ctx context.Context, from *dtc.DtcTopologyRulesInner, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DtcTopologyRulesInnerAttrTypes)
	}
	m := DtcTopologyRulesInnerModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, DtcTopologyRulesInnerAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DtcTopologyRulesInnerModel) Flatten(ctx context.Context, from *dtc.DtcTopologyRulesInner, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DtcTopologyRulesInnerModel{}
	}

	m.DestType = flex.FlattenStringPointer(from.DtcTopologyRulesInnerOneOf1.DestType)
	m.DestinationLink = flex.FlattenStringPointer(from.DtcTopologyRulesInnerOneOf1.DestinationLink)
	m.ReturnType = flex.FlattenStringPointer(from.DtcTopologyRulesInnerOneOf1.ReturnType)
	m.Topology = flex.FlattenStringPointer(from.DtcTopologyRulesInnerOneOf1.Topology)
	m.Valid = types.BoolPointerValue(from.DtcTopologyRulesInnerOneOf1.Valid)
	m.Sources = flex.FlattenFrameworkListNestedBlock(ctx, from.DtcTopologyRulesInnerOneOf1.Sources, DtcTopologyRulesInnerOneOf1SourcesInnerAttrTypes, diags, FlattenDtcTopologyRulesInnerOneOf1SourcesInner)
}

func (m *DtcTopologyRulesInnerModel) PutExpand(to *dtc.DtcTopologyRulesInner) *dtc.DtcTopologyRulesInner {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()

	// Helper to recursively delete empty fields in structs
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

	for field, attr := range DtcTopologyRulesInnerResourceSchemaAttributes {
		attrVal := reflect.ValueOf(attr)
		attrType := attrVal.Type()
		if toType.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < toType.NumField(); i++ {
			tField := toType.Field(i)
			fieldValue := toVal.Field(i).Interface()
			cleanTag := strings.Split(tField.Tag.Get("json"), ",")[0]
			cleanTag = strings.Trim(cleanTag, "_")
			txtFieldValue := utils.ToString(field, fieldValue)
			if field != cleanTag {
				continue
			}

			// Skip if attribute is Required
			if _, ok := attrType.FieldByName("Required"); ok {
				requiredVal := attrVal.FieldByName("Required")
				if requiredVal.IsValid() && requiredVal.CanInterface() {
					boolReq, ok := requiredVal.Interface().(bool)
					if ok && boolReq {
						continue
					}
				}
			}

			// Handle Default
			if _, ok := attrType.FieldByName("Default"); ok {
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

			// Handle Computed
			if _, ok := attrType.FieldByName("Computed"); ok {
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

			// Recursively clean up nested structs and slices
			fvType := reflect.TypeOf(fieldValue)
			if fvType != nil {
				switch fvType.Kind() {
				case reflect.Struct:
					deleteEmptyFields(reflect.ValueOf(fieldValue))
				case reflect.Slice, reflect.Array:
					sliceVal := reflect.ValueOf(fieldValue)
					for j := 0; j < sliceVal.Len(); j++ {
						elem := sliceVal.Index(j)
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
	return to
}
