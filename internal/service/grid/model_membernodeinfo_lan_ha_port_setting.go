package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MembernodeinfoLanHaPortSettingModel struct {
	MgmtLan          iptypes.IPv4Address `tfsdk:"mgmt_lan"`
	MgmtIpv6addr     iptypes.IPv6Address `tfsdk:"mgmt_ipv6addr"`
	HaIpAddress      iptypes.IPAddress   `tfsdk:"ha_ip_address"`
	LanPortSetting   types.Object        `tfsdk:"lan_port_setting"`
	HaPortSetting    types.Object        `tfsdk:"ha_port_setting"`
	HaCloudAttribute types.String        `tfsdk:"ha_cloud_attribute"`
}

var MembernodeinfoLanHaPortSettingAttrTypes = map[string]attr.Type{
	"mgmt_lan":           iptypes.IPv4AddressType{},
	"mgmt_ipv6addr":      iptypes.IPv6AddressType{},
	"ha_ip_address":      iptypes.IPAddressType{},
	"lan_port_setting":   types.ObjectType{AttrTypes: MembernodeinfolanhaportsettingLanPortSettingAttrTypes},
	"ha_port_setting":    types.ObjectType{AttrTypes: MembernodeinfolanhaportsettingHaPortSettingAttrTypes},
	"ha_cloud_attribute": types.StringType,
}

var MembernodeinfoLanHaPortSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"mgmt_lan": schema.StringAttribute{
		CustomType:          iptypes.IPv4AddressType{},
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "Public IPv4 address for the LAN1 interface.",
	},
	"mgmt_ipv6addr": schema.StringAttribute{
		CustomType:          iptypes.IPv6AddressType{},
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "Public IPv6 address for the LAN1 interface.",
	},
	"ha_ip_address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "HA IP address.",
	},
	"lan_port_setting": schema.SingleNestedAttribute{
		Attributes:          MembernodeinfolanhaportsettingLanPortSettingResourceSchemaAttributes,
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Physical port settings for the LAN interface.",
	},
	"ha_port_setting": schema.SingleNestedAttribute{
		Attributes:          MembernodeinfolanhaportsettingHaPortSettingResourceSchemaAttributes,
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Physical port settings for the HA interface.",
	},
	"ha_cloud_attribute": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString("UNK"),
		MarkdownDescription: "HA cloud interface from cloud platform side.",
	},
}

func ExpandMembernodeinfoLanHaPortSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MembernodeinfoLanHaPortSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MembernodeinfoLanHaPortSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MembernodeinfoLanHaPortSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MembernodeinfoLanHaPortSetting {
	if m == nil {
		return nil
	}
	to := &grid.MembernodeinfoLanHaPortSetting{
		MgmtLan:          flex.ExpandIPv4Address(m.MgmtLan),
		MgmtIpv6addr:     flex.ExpandIPv6Address(m.MgmtIpv6addr),
		HaIpAddress:      flex.ExpandIPAddress(m.HaIpAddress),
		LanPortSetting:   ExpandMembernodeinfolanhaportsettingLanPortSetting(ctx, m.LanPortSetting, diags),
		HaPortSetting:    ExpandMembernodeinfolanhaportsettingHaPortSetting(ctx, m.HaPortSetting, diags),
		HaCloudAttribute: flex.ExpandStringPointer(m.HaCloudAttribute),
	}
	return to
}

func FlattenMembernodeinfoLanHaPortSetting(ctx context.Context, from *grid.MembernodeinfoLanHaPortSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MembernodeinfoLanHaPortSettingAttrTypes)
	}
	m := MembernodeinfoLanHaPortSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MembernodeinfoLanHaPortSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MembernodeinfoLanHaPortSettingModel) Flatten(ctx context.Context, from *grid.MembernodeinfoLanHaPortSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MembernodeinfoLanHaPortSettingModel{}
	}
	m.MgmtLan = flex.FlattenIPv4Address(from.MgmtLan)
	m.MgmtIpv6addr = flex.FlattenIPv6Address(from.MgmtIpv6addr)
	m.HaIpAddress = flex.FlattenIPAddress(from.HaIpAddress)
	m.LanPortSetting = FlattenMembernodeinfolanhaportsettingLanPortSetting(ctx, from.LanPortSetting, diags)
	m.HaPortSetting = FlattenMembernodeinfolanhaportsettingHaPortSetting(ctx, from.HaPortSetting, diags)
	m.HaCloudAttribute = flex.FlattenStringPointer(from.HaCloudAttribute)
}

func (m *MembernodeinfoLanHaPortSettingModel) PutExpand(to *grid.MembernodeinfoLanHaPortSetting) *grid.MembernodeinfoLanHaPortSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()

	// Helper to recursively delete empty fields in structs
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

	for field, attr := range MembernodeinfoLanHaPortSettingResourceSchemaAttributes {
		attrVal := reflect.ValueOf(attr)
		attrType := attrVal.Type()
		if toType.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < toType.NumField(); i++ {
			tField := toType.Field(i)
			fieldValue := toVal.Field(i).Interface()
			cleanTag := strings.Split(tField.Tag.Get("json"), ",")[0]
			cleanTag = strings.Trim(cleanTag, "_")
			txtFieldValue := utils.ToString(field, fieldValue)
			if field != cleanTag {
				continue
			}

			// Skip if attribute is Required
			if _, ok := attrType.FieldByName("Required"); ok {
				requiredVal := attrVal.FieldByName("Required")
				if requiredVal.IsValid() && requiredVal.CanInterface() {
					boolReq, ok := requiredVal.Interface().(bool)
					if ok && boolReq {
						continue
					}
				}
			}

			// Handle Default
			if _, ok := attrType.FieldByName("Default"); ok {
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

			// Handle Computed
			if _, ok := attrType.FieldByName("Computed"); ok {
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

			// Recursively clean up nested structs and slices
			fvType := reflect.TypeOf(fieldValue)
			if fvType != nil {
				switch fvType.Kind() {
				case reflect.Struct:
					deleteEmptyFields(reflect.ValueOf(fieldValue))
				case reflect.Slice, reflect.Array:
					sliceVal := reflect.ValueOf(fieldValue)
					for j := 0; j < sliceVal.Len(); j++ {
						elem := sliceVal.Index(j)
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
	return to
}
