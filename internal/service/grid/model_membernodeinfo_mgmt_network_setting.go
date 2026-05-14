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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MembernodeinfoMgmtNetworkSettingModel struct {
	Address       iptypes.IPv4Address `tfsdk:"address"`
	Gateway       types.String        `tfsdk:"gateway"`
	SubnetMask    types.String        `tfsdk:"subnet_mask"`
	VlanId        types.Int64         `tfsdk:"vlan_id"`
	Primary       types.Bool          `tfsdk:"primary"`
	Dscp          types.Int64         `tfsdk:"dscp"`
	LanSubnetMask types.String        `tfsdk:"lan_subnet_mask"`
	LanGateway    types.String        `tfsdk:"lan_gateway"`
	UseDscp       types.Bool          `tfsdk:"use_dscp"`
}

var MembernodeinfoMgmtNetworkSettingAttrTypes = map[string]attr.Type{
	"address":         iptypes.IPv4AddressType{},
	"gateway":         types.StringType,
	"subnet_mask":     types.StringType,
	"vlan_id":         types.Int64Type,
	"primary":         types.BoolType,
	"dscp":            types.Int64Type,
	"lan_subnet_mask": types.StringType,
	"lan_gateway":     types.StringType,
	"use_dscp":        types.BoolType,
}

// Defaults Removed Due to conflicts with Mgmt Network Settings
var MembernodeinfoMgmtNetworkSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		CustomType:          iptypes.IPv4AddressType{},
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The IPv4 Address of the Grid Member.",
	},
	"gateway": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The default gateway for the Grid Member.",
	},
	"subnet_mask": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The subnet mask for the Grid Member.",
	},
	"vlan_id": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The identifier for the VLAN. Valid values are from 1 to 4096.",
	},
	"primary": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if the current address is the primary VLAN address or not.",
	},
	"dscp": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("use_dscp")),
		},
		MarkdownDescription: "The DSCP (Differentiated Services Code Point) value determines relative priorities for the type of services on your network. The appliance implements QoS (Quality of Service) rules based on this configuration. Valid values are from 0 to 63.",
	},
	"lan_subnet_mask": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "LAN netmask only for GCP HA.",
	},
	"lan_gateway": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "LAN gateway only for GCP HA.",
	},
	"use_dscp": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Use flag for: dscp",
	},
}

func ExpandMembernodeinfoMgmtNetworkSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MembernodeinfoMgmtNetworkSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MembernodeinfoMgmtNetworkSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MembernodeinfoMgmtNetworkSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MembernodeinfoMgmtNetworkSetting {
	if m == nil {
		return nil
	}
	to := &grid.MembernodeinfoMgmtNetworkSetting{
		Address:       flex.ExpandIPv4Address(m.Address),
		Gateway:       flex.ExpandStringPointer(m.Gateway),
		SubnetMask:    flex.ExpandStringPointer(m.SubnetMask),
		VlanId:        flex.ExpandInt64Pointer(m.VlanId),
		Primary:       flex.ExpandBoolPointer(m.Primary),
		Dscp:          flex.ExpandInt64Pointer(m.Dscp),
		LanSubnetMask: flex.ExpandStringPointer(m.LanSubnetMask),
		LanGateway:    flex.ExpandStringPointer(m.LanGateway),
		UseDscp:       flex.ExpandBoolPointer(m.UseDscp),
	}
	return to
}

func FlattenMembernodeinfoMgmtNetworkSetting(ctx context.Context, from *grid.MembernodeinfoMgmtNetworkSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MembernodeinfoMgmtNetworkSettingAttrTypes)
	}
	m := MembernodeinfoMgmtNetworkSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MembernodeinfoMgmtNetworkSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MembernodeinfoMgmtNetworkSettingModel) Flatten(ctx context.Context, from *grid.MembernodeinfoMgmtNetworkSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MembernodeinfoMgmtNetworkSettingModel{}
	}
	m.Address = flex.FlattenIPv4Address(from.Address)
	m.Gateway = flex.FlattenStringPointer(from.Gateway)
	m.SubnetMask = flex.FlattenStringPointer(from.SubnetMask)
	m.VlanId = flex.FlattenInt64Pointer(from.VlanId)
	m.Primary = types.BoolPointerValue(from.Primary)
	m.Dscp = flex.FlattenInt64Pointer(from.Dscp)
	m.LanSubnetMask = flex.FlattenStringPointer(from.LanSubnetMask)
	m.LanGateway = flex.FlattenStringPointer(from.LanGateway)
	m.UseDscp = types.BoolPointerValue(from.UseDscp)
}

func (m *MembernodeinfoMgmtNetworkSettingModel) PutExpand(to *grid.MembernodeinfoMgmtNetworkSetting) *grid.MembernodeinfoMgmtNetworkSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MembernodeinfoMgmtNetworkSettingResourceSchemaAttributes {
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
