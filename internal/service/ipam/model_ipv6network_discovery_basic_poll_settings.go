package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type Ipv6networkDiscoveryBasicPollSettingsModel struct {
	PortScanning                            types.Bool   `tfsdk:"port_scanning"`
	DeviceProfile                           types.Bool   `tfsdk:"device_profile"`
	SnmpCollection                          types.Bool   `tfsdk:"snmp_collection"`
	CliCollection                           types.Bool   `tfsdk:"cli_collection"`
	NetbiosScanning                         types.Bool   `tfsdk:"netbios_scanning"`
	CompletePingSweep                       types.Bool   `tfsdk:"complete_ping_sweep"`
	SmartSubnetPingSweep                    types.Bool   `tfsdk:"smart_subnet_ping_sweep"`
	AutoArpRefreshBeforeSwitchPortPolling   types.Bool   `tfsdk:"auto_arp_refresh_before_switch_port_polling"`
	SwitchPortDataCollectionPolling         types.String `tfsdk:"switch_port_data_collection_polling"`
	SwitchPortDataCollectionPollingSchedule types.Object `tfsdk:"switch_port_data_collection_polling_schedule"`
	SwitchPortDataCollectionPollingInterval types.Int64  `tfsdk:"switch_port_data_collection_polling_interval"`
	CredentialGroup                         types.String `tfsdk:"credential_group"`
	PollingFrequencyModifier                types.String `tfsdk:"polling_frequency_modifier"`
	UseGlobalPollingFrequencyModifier       types.Bool   `tfsdk:"use_global_polling_frequency_modifier"`
}

var Ipv6networkDiscoveryBasicPollSettingsAttrTypes = map[string]attr.Type{
	"port_scanning":                                types.BoolType,
	"device_profile":                               types.BoolType,
	"snmp_collection":                              types.BoolType,
	"cli_collection":                               types.BoolType,
	"netbios_scanning":                             types.BoolType,
	"complete_ping_sweep":                          types.BoolType,
	"smart_subnet_ping_sweep":                      types.BoolType,
	"auto_arp_refresh_before_switch_port_polling":  types.BoolType,
	"switch_port_data_collection_polling":          types.StringType,
	"switch_port_data_collection_polling_schedule": types.ObjectType{AttrTypes: Ipv6networkdiscoverybasicpollsettingsSwitchPortDataCollectionPollingScheduleAttrTypes},
	"switch_port_data_collection_polling_interval": types.Int64Type,
	"credential_group":                             types.StringType,
	"polling_frequency_modifier":                   types.StringType,
	"use_global_polling_frequency_modifier":        types.BoolType,
}

var Ipv6networkDiscoveryBasicPollSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"port_scanning": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether port scanning is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"device_profile": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether device profile is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"snmp_collection": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether SNMP collection is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
	"cli_collection": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether CLI collection is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
	"netbios_scanning": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether netbios scanning is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"complete_ping_sweep": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether complete ping sweep is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"smart_subnet_ping_sweep": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether smart subnet ping sweep is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"auto_arp_refresh_before_switch_port_polling": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether auto ARP refresh before switch port polling is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
	"switch_port_data_collection_polling": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "A switch port data collection polling mode.",
		Computed:            true,
		Default:             stringdefault.StaticString("PERIODIC"),
		Validators: []validator.String{
			stringvalidator.OneOf(
				"PERIODIC",
				"DISABLED",
				"SCHEDULED",
			),
		},
	},
	"switch_port_data_collection_polling_schedule": schema.SingleNestedAttribute{
		Attributes:          Ipv6networkdiscoverybasicpollsettingsSwitchPortDataCollectionPollingScheduleResourceSchemaAttributes,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "A switch port data collection polling schedule.",
	},
	"switch_port_data_collection_polling_interval": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Indicates the interval for switch port data collection polling.",
		Computed:            true,
		Default:             int64default.StaticInt64(3600),
	},
	"credential_group": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Credential group.",
		Computed:            true,
		Default:             stringdefault.StaticString("default"),
	},
	"polling_frequency_modifier": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Polling Frequency Modifier.",
		Computed:            true,
		Default:             stringdefault.StaticString("1"),
	},
	"use_global_polling_frequency_modifier": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Use Global Polling Frequency Modifier.",
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
}

func ExpandIpv6networkDiscoveryBasicPollSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.Ipv6networkDiscoveryBasicPollSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m Ipv6networkDiscoveryBasicPollSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *Ipv6networkDiscoveryBasicPollSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.Ipv6networkDiscoveryBasicPollSettings {
	if m == nil {
		return nil
	}
	to := &ipam.Ipv6networkDiscoveryBasicPollSettings{
		PortScanning:                            flex.ExpandBoolPointer(m.PortScanning),
		DeviceProfile:                           flex.ExpandBoolPointer(m.DeviceProfile),
		SnmpCollection:                          flex.ExpandBoolPointer(m.SnmpCollection),
		CliCollection:                           flex.ExpandBoolPointer(m.CliCollection),
		NetbiosScanning:                         flex.ExpandBoolPointer(m.NetbiosScanning),
		CompletePingSweep:                       flex.ExpandBoolPointer(m.CompletePingSweep),
		SmartSubnetPingSweep:                    flex.ExpandBoolPointer(m.SmartSubnetPingSweep),
		AutoArpRefreshBeforeSwitchPortPolling:   flex.ExpandBoolPointer(m.AutoArpRefreshBeforeSwitchPortPolling),
		SwitchPortDataCollectionPolling:         flex.ExpandStringPointer(m.SwitchPortDataCollectionPolling),
		SwitchPortDataCollectionPollingSchedule: ExpandIpv6networkdiscoverybasicpollsettingsSwitchPortDataCollectionPollingSchedule(ctx, m.SwitchPortDataCollectionPollingSchedule, diags),
		SwitchPortDataCollectionPollingInterval: flex.ExpandInt64Pointer(m.SwitchPortDataCollectionPollingInterval),
		CredentialGroup:                         flex.ExpandStringPointer(m.CredentialGroup),
		PollingFrequencyModifier:                flex.ExpandStringPointer(m.PollingFrequencyModifier),
		UseGlobalPollingFrequencyModifier:       flex.ExpandBoolPointer(m.UseGlobalPollingFrequencyModifier),
	}
	return to
}

func FlattenIpv6networkDiscoveryBasicPollSettings(ctx context.Context, from *ipam.Ipv6networkDiscoveryBasicPollSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(Ipv6networkDiscoveryBasicPollSettingsAttrTypes)
	}
	m := Ipv6networkDiscoveryBasicPollSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, Ipv6networkDiscoveryBasicPollSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *Ipv6networkDiscoveryBasicPollSettingsModel) Flatten(ctx context.Context, from *ipam.Ipv6networkDiscoveryBasicPollSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = Ipv6networkDiscoveryBasicPollSettingsModel{}
	}
	m.PortScanning = types.BoolPointerValue(from.PortScanning)
	m.DeviceProfile = types.BoolPointerValue(from.DeviceProfile)
	m.SnmpCollection = types.BoolPointerValue(from.SnmpCollection)
	m.CliCollection = types.BoolPointerValue(from.CliCollection)
	m.NetbiosScanning = types.BoolPointerValue(from.NetbiosScanning)
	m.CompletePingSweep = types.BoolPointerValue(from.CompletePingSweep)
	m.SmartSubnetPingSweep = types.BoolPointerValue(from.SmartSubnetPingSweep)
	m.AutoArpRefreshBeforeSwitchPortPolling = types.BoolPointerValue(from.AutoArpRefreshBeforeSwitchPortPolling)
	m.SwitchPortDataCollectionPolling = flex.FlattenStringPointer(from.SwitchPortDataCollectionPolling)
	m.SwitchPortDataCollectionPollingSchedule = FlattenIpv6networkdiscoverybasicpollsettingsSwitchPortDataCollectionPollingSchedule(ctx, from.SwitchPortDataCollectionPollingSchedule, diags)
	m.SwitchPortDataCollectionPollingInterval = flex.FlattenInt64Pointer(from.SwitchPortDataCollectionPollingInterval)
	m.CredentialGroup = flex.FlattenStringPointer(from.CredentialGroup)
	m.PollingFrequencyModifier = flex.FlattenStringPointer(from.PollingFrequencyModifier)
	m.UseGlobalPollingFrequencyModifier = types.BoolPointerValue(from.UseGlobalPollingFrequencyModifier)
}

func (m *Ipv6networkDiscoveryBasicPollSettingsModel) PutExpand(to *ipam.Ipv6networkDiscoveryBasicPollSettings) *ipam.Ipv6networkDiscoveryBasicPollSettings {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range Ipv6networkDiscoveryBasicPollSettingsResourceSchemaAttributes {
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
