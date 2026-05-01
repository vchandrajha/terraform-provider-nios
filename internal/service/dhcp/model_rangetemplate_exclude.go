package dhcp

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type RangetemplateExcludeModel struct {
	Offset            types.Int64  `tfsdk:"offset"`
	NumberOfAddresses types.Int64  `tfsdk:"number_of_addresses"`
	Comment           types.String `tfsdk:"comment"`
}

var RangetemplateExcludeAttrTypes = map[string]attr.Type{
	"offset":              types.Int64Type,
	"number_of_addresses": types.Int64Type,
	"comment":             types.StringType,
}

var RangetemplateExcludeResourceSchemaAttributes = map[string]schema.Attribute{
	"offset": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The address offset of the DHCP exclusion range template.",
	},
	"number_of_addresses": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The number of addresses in the DHCP exclusion range template.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "A descriptive comment of a DHCP exclusion range template.",
	},
}

func ExpandRangetemplateExclude(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dhcp.RangetemplateExclude {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m RangetemplateExcludeModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *RangetemplateExcludeModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.RangetemplateExclude {
	if m == nil {
		return nil
	}
	to := &dhcp.RangetemplateExclude{
		Offset:            flex.ExpandInt64Pointer(m.Offset),
		NumberOfAddresses: flex.ExpandInt64Pointer(m.NumberOfAddresses),
		Comment:           flex.ExpandStringPointer(m.Comment),
	}
	return to
}

func FlattenRangetemplateExclude(ctx context.Context, from *dhcp.RangetemplateExclude, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RangetemplateExcludeAttrTypes)
	}
	m := RangetemplateExcludeModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RangetemplateExcludeAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RangetemplateExcludeModel) Flatten(ctx context.Context, from *dhcp.RangetemplateExclude, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RangetemplateExcludeModel{}
	}
	m.Offset = flex.FlattenInt64Pointer(from.Offset)
	m.NumberOfAddresses = flex.FlattenInt64Pointer(from.NumberOfAddresses)
	m.Comment = flex.FlattenStringPointer(from.Comment)
}

func (m *RangetemplateExcludeModel) PutExpand(to *dhcp.RangetemplateExclude) *dhcp.RangetemplateExclude {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RangetemplateExcludeResourceSchemaAttributes {
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
