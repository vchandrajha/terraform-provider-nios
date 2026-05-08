package dhcp

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type Ipv6dhcpoptiondefinitionModel struct {
	Ref   types.String `tfsdk:"ref"`
	Code  types.Int64  `tfsdk:"code"`
	Name  types.String `tfsdk:"name"`
	Space types.String `tfsdk:"space"`
	Type  types.String `tfsdk:"type"`
}

var Ipv6dhcpoptiondefinitionAttrTypes = map[string]attr.Type{
	"ref":   types.StringType,
	"code":  types.Int64Type,
	"name":  types.StringType,
	"space": types.StringType,
	"type":  types.StringType,
}

var Ipv6dhcpoptiondefinitionResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"code": schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
		MarkdownDescription: "The code of a DHCP IPv6 option definition object. An option code number is used to identify the DHCP option.",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name of a DHCP IPv6 option definition object.",
	},
	"space": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The space of a DHCP option definition object.",
	},
	"type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("16-bit signed integer", "16-bit unsigned integer", "32-bit signed integer",
				"32-bit unsigned integer", "8-bit signed integer", "8-bit unsigned integer", "8-bit unsigned integer",
				"array of 16-bit integer", "array of 16-bit unsigned integer", "array of 32-bit integer",
				"array of 32-bit unsigned integer", "array of 8-bit integer", "array of 8-bit unsigned integer",
				"array of ip-address", "boolean", "boolean array of ip-address", "boolean-text",
				"domain-list", "domain-name", "ip-address", "string", "text",
			),
		},
		MarkdownDescription: "The data type of the Grid DHCP IPv6 option.",
	},
}

func (m *Ipv6dhcpoptiondefinitionModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.Ipv6dhcpoptiondefinition {
	if m == nil {
		return nil
	}
	to := &dhcp.Ipv6dhcpoptiondefinition{
		Code:  flex.ExpandInt64Pointer(m.Code),
		Name:  flex.ExpandStringPointer(m.Name),
		Space: flex.ExpandStringPointer(m.Space),
		Type:  flex.ExpandStringPointer(m.Type),
	}
	return to
}

func FlattenIpv6dhcpoptiondefinition(ctx context.Context, from *dhcp.Ipv6dhcpoptiondefinition, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(Ipv6dhcpoptiondefinitionAttrTypes)
	}
	m := Ipv6dhcpoptiondefinitionModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, Ipv6dhcpoptiondefinitionAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *Ipv6dhcpoptiondefinitionModel) Flatten(ctx context.Context, from *dhcp.Ipv6dhcpoptiondefinition, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = Ipv6dhcpoptiondefinitionModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Code = flex.FlattenInt64Pointer(from.Code)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Space = flex.FlattenStringPointer(from.Space)
	m.Type = flex.FlattenStringPointer(from.Type)
}

func (m *Ipv6dhcpoptiondefinitionModel) PutExpand(to *dhcp.Ipv6dhcpoptiondefinition) *dhcp.Ipv6dhcpoptiondefinition {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range Ipv6dhcpoptiondefinitionResourceSchemaAttributes {
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
