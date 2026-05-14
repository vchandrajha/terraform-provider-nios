package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberLan2PortSettingModel struct {
	VirtualRouterId             types.Int64  `tfsdk:"virtual_router_id"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	NetworkSetting              types.Object `tfsdk:"network_setting"`
	V6NetworkSetting            types.Object `tfsdk:"v6_network_setting"`
	NicFailoverEnabled          types.Bool   `tfsdk:"nic_failover_enabled"`
	NicFailoverEnablePrimary    types.Bool   `tfsdk:"nic_failover_enable_primary"`
	DefaultRouteFailoverEnabled types.Bool   `tfsdk:"default_route_failover_enabled"`
}

var MemberLan2PortSettingAttrTypes = map[string]attr.Type{
	"virtual_router_id":              types.Int64Type,
	"enabled":                        types.BoolType,
	"network_setting":                types.ObjectType{AttrTypes: Memberlan2portsettingNetworkSettingAttrTypes},
	"v6_network_setting":             types.ObjectType{AttrTypes: Memberlan2portsettingV6NetworkSettingAttrTypes},
	"nic_failover_enabled":           types.BoolType,
	"nic_failover_enable_primary":    types.BoolType,
	"default_route_failover_enabled": types.BoolType,
}

var MemberLan2PortSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"virtual_router_id": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If the 'enabled' field is set to True, this defines the virtual router ID for the LAN2 port.",
	},
	"enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "If this field is set to True, then it has its own IP settings. Otherwise, port redundancy mechanism is used, in which the LAN1 and LAN2 ports share the same IP settings for failover purposes.",
	},
	"network_setting": schema.SingleNestedAttribute{
		Attributes:          Memberlan2portsettingNetworkSettingResourceSchemaAttributes,
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "If the ‘enable’ field is set to True, this defines IPv4 network settings for LAN2.",
	},
	"v6_network_setting": schema.SingleNestedAttribute{
		Attributes:          Memberlan2portsettingV6NetworkSettingResourceSchemaAttributes,
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "If the ‘enable’ field is set to True, this defines IPv6 network settings for the LAN2 port.",
	},
	"nic_failover_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if NIC failover is enabled or not.",
	},
	"nic_failover_enable_primary": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Prefer LAN1 when available.",
	},
	"default_route_failover_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Default route failover for LAN1 and LAN2.",
	},
}

func ExpandMemberLan2PortSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberLan2PortSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberLan2PortSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberLan2PortSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberLan2PortSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberLan2PortSetting{
		VirtualRouterId:             flex.ExpandInt64Pointer(m.VirtualRouterId),
		Enabled:                     flex.ExpandBoolPointer(m.Enabled),
		NetworkSetting:              ExpandMemberlan2portsettingNetworkSetting(ctx, m.NetworkSetting, diags),
		V6NetworkSetting:            ExpandMemberlan2portsettingV6NetworkSetting(ctx, m.V6NetworkSetting, diags),
		NicFailoverEnabled:          flex.ExpandBoolPointer(m.NicFailoverEnabled),
		NicFailoverEnablePrimary:    flex.ExpandBoolPointer(m.NicFailoverEnablePrimary),
		DefaultRouteFailoverEnabled: flex.ExpandBoolPointer(m.DefaultRouteFailoverEnabled),
	}
	return to
}

func FlattenMemberLan2PortSetting(ctx context.Context, from *grid.MemberLan2PortSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberLan2PortSettingAttrTypes)
	}
	m := MemberLan2PortSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberLan2PortSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberLan2PortSettingModel) Flatten(ctx context.Context, from *grid.MemberLan2PortSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberLan2PortSettingModel{}
	}
	m.VirtualRouterId = flex.FlattenInt64Pointer(from.VirtualRouterId)
	m.Enabled = types.BoolPointerValue(from.Enabled)
	m.NetworkSetting = FlattenMemberlan2portsettingNetworkSetting(ctx, from.NetworkSetting, diags)
	m.V6NetworkSetting = FlattenMemberlan2portsettingV6NetworkSetting(ctx, from.V6NetworkSetting, diags)
	m.NicFailoverEnabled = types.BoolPointerValue(from.NicFailoverEnabled)
	m.NicFailoverEnablePrimary = types.BoolPointerValue(from.NicFailoverEnablePrimary)
	m.DefaultRouteFailoverEnabled = types.BoolPointerValue(from.DefaultRouteFailoverEnabled)
}

func (m *MemberLan2PortSettingModel) PutExpand(to *grid.MemberLan2PortSetting) *grid.MemberLan2PortSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberLan2PortSettingResourceSchemaAttributes {
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
