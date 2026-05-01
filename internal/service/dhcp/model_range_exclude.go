package dhcp

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type RangeExcludeModel struct {
	StartAddress iptypes.IPv4Address `tfsdk:"start_address"`
	EndAddress   iptypes.IPv4Address `tfsdk:"end_address"`
	Comment      types.String        `tfsdk:"comment"`
}

var RangeExcludeAttrTypes = map[string]attr.Type{
	"start_address": iptypes.IPv4AddressType{},
	"end_address":   iptypes.IPv4AddressType{},
	"comment":       types.StringType,
}

var RangeExcludeResourceSchemaAttributes = map[string]schema.Attribute{
	"start_address": schema.StringAttribute{
		CustomType:          iptypes.IPv4AddressType{},
		Required:            true,
		MarkdownDescription: "The IPv4 Address starting address of the exclusion range.",
	},
	"end_address": schema.StringAttribute{
		CustomType:          iptypes.IPv4AddressType{},
		Required:            true,
		MarkdownDescription: "The IPv4 Address ending address of the exclusion range.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Comment for the exclusion range; maximum 256 characters.",
	},
}

func ExpandRangeExclude(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dhcp.RangeExclude {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m RangeExcludeModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *RangeExcludeModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.RangeExclude {
	if m == nil {
		return nil
	}
	to := &dhcp.RangeExclude{
		StartAddress: flex.ExpandIPv4Address(m.StartAddress),
		EndAddress:   flex.ExpandIPv4Address(m.EndAddress),
		Comment:      flex.ExpandStringPointer(m.Comment),
	}
	return to
}

func FlattenRangeExclude(ctx context.Context, from *dhcp.RangeExclude, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RangeExcludeAttrTypes)
	}
	m := RangeExcludeModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RangeExcludeAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RangeExcludeModel) Flatten(ctx context.Context, from *dhcp.RangeExclude, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RangeExcludeModel{}
	}
	m.StartAddress = flex.FlattenIPv4Address(from.StartAddress)
	m.EndAddress = flex.FlattenIPv4Address(from.EndAddress)
	m.Comment = flex.FlattenStringPointer(from.Comment)
}

func (m *RangeExcludeModel) PutExpand(to *dhcp.RangeExclude) *dhcp.RangeExclude {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RangeExcludeResourceSchemaAttributes {
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
