package dtc

import (
	"context"
	"reflect"
	"strings"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type DtcTopologyRulesInnerOneOf1SourcesInnerModel struct {
	SourceOp    types.String `tfsdk:"source_op"`
	SourceType  types.String `tfsdk:"source_type"`
	SourceValue types.String `tfsdk:"source_value"`
}

var DtcTopologyRulesInnerOneOf1SourcesInnerAttrTypes = map[string]attr.Type{
	"source_op":    types.StringType,
	"source_type":  types.StringType,
	"source_value": types.StringType,
}

var DtcTopologyRulesInnerOneOf1SourcesInnerResourceSchemaAttributes = map[string]schema.Attribute{
	"source_op": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("IS", "IS_NOT"),
		},
		MarkdownDescription: "Operation for matching the source.",
	},
	"source_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("CITY", "CONTINENT", "COUNTRY", "EA0", "EA1", "EA2", "EA3", "SUBDIVISION", "SUBNET"),
		},
		MarkdownDescription: "Type of the source.",
	},
	"source_value": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Value of the source.",
	},
}

func ExpandDtcTopologyRulesInnerOneOf1SourcesInner(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dtc.DtcTopologyRulesInnerOneOf1SourcesInner {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m DtcTopologyRulesInnerOneOf1SourcesInnerModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *DtcTopologyRulesInnerOneOf1SourcesInnerModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dtc.DtcTopologyRulesInnerOneOf1SourcesInner {
	if m == nil {
		return nil
	}
	to := &dtc.DtcTopologyRulesInnerOneOf1SourcesInner{
		SourceOp:    flex.ExpandStringPointer(m.SourceOp),
		SourceType:  flex.ExpandStringPointer(m.SourceType),
		SourceValue: flex.ExpandStringPointer(m.SourceValue),
	}
	return to
}

func FlattenDtcTopologyRulesInnerOneOf1SourcesInner(ctx context.Context, from *dtc.DtcTopologyRulesInnerOneOf1SourcesInner, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DtcTopologyRulesInnerOneOf1SourcesInnerAttrTypes)
	}
	m := DtcTopologyRulesInnerOneOf1SourcesInnerModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, DtcTopologyRulesInnerOneOf1SourcesInnerAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DtcTopologyRulesInnerOneOf1SourcesInnerModel) Flatten(ctx context.Context, from *dtc.DtcTopologyRulesInnerOneOf1SourcesInner, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DtcTopologyRulesInnerOneOf1SourcesInnerModel{}
	}
	m.SourceOp = flex.FlattenStringPointer(from.SourceOp)
	m.SourceType = flex.FlattenStringPointer(from.SourceType)
	m.SourceValue = flex.FlattenStringPointer(from.SourceValue)
}

func (m *DtcTopologyRulesInnerOneOf1SourcesInnerModel) PutExpand(to *dtc.DtcTopologyRulesInnerOneOf1SourcesInner) *dtc.DtcTopologyRulesInnerOneOf1SourcesInner {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range DtcTopologyRulesInnerOneOf1SourcesInnerResourceSchemaAttributes {
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
