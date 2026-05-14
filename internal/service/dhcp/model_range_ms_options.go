package dhcp

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type RangeMsOptionsModel struct {
	Num         types.Int64  `tfsdk:"num"`
	Value       types.String `tfsdk:"value"`
	Name        types.String `tfsdk:"name"`
	VendorClass types.String `tfsdk:"vendor_class"`
	UserClass   types.String `tfsdk:"user_class"`
	Type        types.String `tfsdk:"type"`
}

var RangeMsOptionsAttrTypes = map[string]attr.Type{
	"num":          types.Int64Type,
	"value":        types.StringType,
	"name":         types.StringType,
	"vendor_class": types.StringType,
	"user_class":   types.StringType,
	"type":         types.StringType,
}

var RangeMsOptionsResourceSchemaAttributes = map[string]schema.Attribute{
	"num": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The code of the DHCP option.",
	},
	"value": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Value of the DHCP option. Required to be set for all options.",
	},
	"name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The name of the DHCP option.",
	},
	"vendor_class": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The name of the vendor class with which this DHCP option is associated.",
	},
	"user_class": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The name of the user class with which this DHCP option is associated.",
	},
	"type": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The DHCP option type. Valid values are: * \"16-bit signed integer\" * \"16-bit unsigned integer\" * \"32-bit signed integer\" * \"32-bit unsigned integer\" * \"64-bit unsigned integer\" * \"8-bit signed integer\" * \"8-bit unsigned integer (1,2,4,8)\" * \"8-bit unsigned integer\" * \"array of 16-bit integer\" * \"array of 16-bit unsigned integer\" * \"array of 32-bit integer\" * \"array of 32-bit unsigned integer\" * \"array of 64-bit unsigned integer\" * \"array of 8-bit integer\" * \"array of 8-bit unsigned integer\" * \"array of ip-address pair\" * \"array of ip-address\" * \"array of string\" * \"binary\" * \"boolean array of ip-address\" * \"boolean\" * \"boolean-text\" * \"domain-list\" * \"domain-name\" * \"encapsulated\" * \"ip-address\" * \"string\" * \"text\"",
	},
}

func ExpandRangeMsOptions(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dhcp.RangeMsOptions {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m RangeMsOptionsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *RangeMsOptionsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.RangeMsOptions {
	if m == nil {
		return nil
	}
	to := &dhcp.RangeMsOptions{
		Num:         flex.ExpandInt64Pointer(m.Num),
		Value:       flex.ExpandStringPointer(m.Value),
		Name:        flex.ExpandStringPointer(m.Name),
		VendorClass: flex.ExpandStringPointer(m.VendorClass),
		UserClass:   flex.ExpandStringPointer(m.UserClass),
	}
	return to
}

func FlattenRangeMsOptions(ctx context.Context, from *dhcp.RangeMsOptions, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RangeMsOptionsAttrTypes)
	}
	m := RangeMsOptionsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RangeMsOptionsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RangeMsOptionsModel) Flatten(ctx context.Context, from *dhcp.RangeMsOptions, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RangeMsOptionsModel{}
	}
	m.Num = flex.FlattenInt64Pointer(from.Num)
	m.Value = flex.FlattenStringPointer(from.Value)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.VendorClass = flex.FlattenStringPointer(from.VendorClass)
	m.UserClass = flex.FlattenStringPointer(from.UserClass)
	m.Type = flex.FlattenStringPointer(from.Type)
}

func (m *RangeMsOptionsModel) PutExpand(to *dhcp.RangeMsOptions) *dhcp.RangeMsOptions {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RangeMsOptionsResourceSchemaAttributes {
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
