package ipam

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type NetworkSubscribeSettingsModel struct {
	EnabledAttributes  types.List `tfsdk:"enabled_attributes"`
	MappedEaAttributes types.List `tfsdk:"mapped_ea_attributes"`
}

var NetworkSubscribeSettingsAttrTypes = map[string]attr.Type{
	"enabled_attributes":   types.ListType{ElemType: types.StringType},
	"mapped_ea_attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksubscribesettingsMappedEaAttributesAttrTypes}},
}

var NetworkSubscribeSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"enabled_attributes": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The list of Cisco ISE attributes allowed for subscription.",
		Validators: []validator.List{
			customvalidator.StringsInSlice([]string{
				"DOMAINNAME",
				"ENDPOINT_PROFILE",
				"SECURITY_GROUP",
				"SESSION_STATE",
				"SSID",
				"USERNAME",
				"VLAN",
			}),
			listvalidator.SizeAtLeast(1),
		},
	},
	"mapped_ea_attributes": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: NetworksubscribesettingsMappedEaAttributesResourceSchemaAttributes,
		},
		Optional:            true,
		MarkdownDescription: "The list of NIOS extensible attributes to Cisco ISE attributes mappings.",
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
}

func ExpandNetworkSubscribeSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworkSubscribeSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworkSubscribeSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworkSubscribeSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworkSubscribeSettings {
	if m == nil {
		return nil
	}
	to := &ipam.NetworkSubscribeSettings{
		EnabledAttributes:  flex.ExpandFrameworkListString(ctx, m.EnabledAttributes, diags),
		MappedEaAttributes: flex.ExpandFrameworkListNestedBlock(ctx, m.MappedEaAttributes, diags, ExpandNetworksubscribesettingsMappedEaAttributes),
	}
	return to
}

func FlattenNetworkSubscribeSettings(ctx context.Context, from *ipam.NetworkSubscribeSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworkSubscribeSettingsAttrTypes)
	}
	m := NetworkSubscribeSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworkSubscribeSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworkSubscribeSettingsModel) Flatten(ctx context.Context, from *ipam.NetworkSubscribeSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworkSubscribeSettingsModel{}
	}
	m.EnabledAttributes = flex.FlattenFrameworkListString(ctx, from.EnabledAttributes, diags)
	m.MappedEaAttributes = flex.FlattenFrameworkListNestedBlock(ctx, from.MappedEaAttributes, NetworksubscribesettingsMappedEaAttributesAttrTypes, diags, FlattenNetworksubscribesettingsMappedEaAttributes)
}

func (m *NetworkSubscribeSettingsModel) PutExpand(to *ipam.NetworkSubscribeSettings) *ipam.NetworkSubscribeSettings {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NetworkSubscribeSettingsResourceSchemaAttributes {
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
