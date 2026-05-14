package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type Memberlan2portsettingV6NetworkSettingModel struct {
	Enabled                 types.Bool          `tfsdk:"enabled"`
	VirtualIp               iptypes.IPv6Address `tfsdk:"virtual_ip"`
	CidrPrefix              types.Int64         `tfsdk:"cidr_prefix"`
	Gateway                 types.String        `tfsdk:"gateway"`
	AutoRouterConfigEnabled types.Bool          `tfsdk:"auto_router_config_enabled"`
	VlanId                  types.Int64         `tfsdk:"vlan_id"`
	Primary                 types.Bool          `tfsdk:"primary"`
	Dscp                    types.Int64         `tfsdk:"dscp"`
	UseDscp                 types.Bool          `tfsdk:"use_dscp"`
}

var Memberlan2portsettingV6NetworkSettingAttrTypes = map[string]attr.Type{
	"enabled":                    types.BoolType,
	"virtual_ip":                 iptypes.IPv6AddressType{},
	"cidr_prefix":                types.Int64Type,
	"gateway":                    types.StringType,
	"auto_router_config_enabled": types.BoolType,
	"vlan_id":                    types.Int64Type,
	"primary":                    types.BoolType,
	"dscp":                       types.Int64Type,
	"use_dscp":                   types.BoolType,
}

var Memberlan2portsettingV6NetworkSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Determines if IPv6 networking should be enabled.",
	},
	"virtual_ip": schema.StringAttribute{
		CustomType:          iptypes.IPv6AddressType{},
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "IPv6 address.",
	},
	"cidr_prefix": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "IPv6 cidr prefix",
	},
	"gateway": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Gateway address.",
	},
	"auto_router_config_enabled": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines if automatic router configuration should be enabled.",
	},
	"vlan_id": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "The identifier for the VLAN. Valid values are from 1 to 4096.",
	},
	"primary": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "Determines if the current address is the primary VLAN address or not.",
	},
	"dscp": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("use_dscp")),
		},
		MarkdownDescription: "The DSCP (Differentiated Services Code Point) value determines relative priorities for the type of services on your network. The appliance implements QoS (Quality of Service) rules based on this configuration. Valid values are from 0 to 63.",
	},
	"use_dscp": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Use flag for: dscp",
	},
}

func ExpandMemberlan2portsettingV6NetworkSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.Memberlan2portsettingV6NetworkSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m Memberlan2portsettingV6NetworkSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *Memberlan2portsettingV6NetworkSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.Memberlan2portsettingV6NetworkSetting {
	if m == nil {
		return nil
	}
	to := &grid.Memberlan2portsettingV6NetworkSetting{
		Enabled:                 flex.ExpandBoolPointer(m.Enabled),
		VirtualIp:               flex.ExpandIPv6Address(m.VirtualIp),
		CidrPrefix:              flex.ExpandInt64Pointer(m.CidrPrefix),
		Gateway:                 flex.ExpandStringPointer(m.Gateway),
		AutoRouterConfigEnabled: flex.ExpandBoolPointer(m.AutoRouterConfigEnabled),
		VlanId:                  flex.ExpandInt64Pointer(m.VlanId),
		Primary:                 flex.ExpandBoolPointer(m.Primary),
		Dscp:                    flex.ExpandInt64Pointer(m.Dscp),
		UseDscp:                 flex.ExpandBoolPointer(m.UseDscp),
	}
	return to
}

func FlattenMemberlan2portsettingV6NetworkSetting(ctx context.Context, from *grid.Memberlan2portsettingV6NetworkSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(Memberlan2portsettingV6NetworkSettingAttrTypes)
	}
	m := Memberlan2portsettingV6NetworkSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, Memberlan2portsettingV6NetworkSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *Memberlan2portsettingV6NetworkSettingModel) Flatten(ctx context.Context, from *grid.Memberlan2portsettingV6NetworkSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = Memberlan2portsettingV6NetworkSettingModel{}
	}
	m.Enabled = types.BoolPointerValue(from.Enabled)
	m.VirtualIp = flex.FlattenIPv6Address(from.VirtualIp)
	m.CidrPrefix = flex.FlattenInt64Pointer(from.CidrPrefix)
	m.Gateway = flex.FlattenStringPointer(from.Gateway)
	m.AutoRouterConfigEnabled = types.BoolPointerValue(from.AutoRouterConfigEnabled)
	m.VlanId = flex.FlattenInt64Pointer(from.VlanId)
	m.Primary = types.BoolPointerValue(from.Primary)
	m.Dscp = flex.FlattenInt64Pointer(from.Dscp)
	m.UseDscp = types.BoolPointerValue(from.UseDscp)
}

func (m *Memberlan2portsettingV6NetworkSettingModel) PutExpand(to *grid.Memberlan2portsettingV6NetworkSetting) *grid.Memberlan2portsettingV6NetworkSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range Memberlan2portsettingV6NetworkSettingResourceSchemaAttributes {
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
