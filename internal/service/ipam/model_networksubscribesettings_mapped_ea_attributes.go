package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NetworksubscribesettingsMappedEaAttributesModel struct {
	Name     types.String `tfsdk:"name"`
	MappedEa types.String `tfsdk:"mapped_ea"`
}

var NetworksubscribesettingsMappedEaAttributesAttrTypes = map[string]attr.Type{
	"name":      types.StringType,
	"mapped_ea": types.StringType,
}

var NetworksubscribesettingsMappedEaAttributesResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The Cisco ISE attribute name that is enabled for publishsing from a Cisco ISE endpoint.",
	},
	"mapped_ea": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The name of the extensible attribute definition object the Cisco ISE attribute that is enabled for subscription is mapped on.",
	},
}

func ExpandNetworksubscribesettingsMappedEaAttributes(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworksubscribesettingsMappedEaAttributes {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworksubscribesettingsMappedEaAttributesModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworksubscribesettingsMappedEaAttributesModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworksubscribesettingsMappedEaAttributes {
	if m == nil {
		return nil
	}
	to := &ipam.NetworksubscribesettingsMappedEaAttributes{
		Name:     flex.ExpandStringPointer(m.Name),
		MappedEa: flex.ExpandStringPointer(m.MappedEa),
	}
	return to
}

func FlattenNetworksubscribesettingsMappedEaAttributes(ctx context.Context, from *ipam.NetworksubscribesettingsMappedEaAttributes, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworksubscribesettingsMappedEaAttributesAttrTypes)
	}
	m := NetworksubscribesettingsMappedEaAttributesModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworksubscribesettingsMappedEaAttributesAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworksubscribesettingsMappedEaAttributesModel) Flatten(ctx context.Context, from *ipam.NetworksubscribesettingsMappedEaAttributes, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworksubscribesettingsMappedEaAttributesModel{}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.MappedEa = flex.FlattenStringPointer(from.MappedEa)
}

func (m *NetworksubscribesettingsMappedEaAttributesModel) PutExpand(to *ipam.NetworksubscribesettingsMappedEaAttributes) *ipam.NetworksubscribesettingsMappedEaAttributes {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NetworksubscribesettingsMappedEaAttributesResourceSchemaAttributes {
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
