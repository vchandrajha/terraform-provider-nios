package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberTrafficCaptureRecDnsSettingModel struct {
	RecDnsLatencyTriggerEnable  types.Bool   `tfsdk:"rec_dns_latency_trigger_enable"`
	RecDnsLatencyThreshold      types.Int64  `tfsdk:"rec_dns_latency_threshold"`
	RecDnsLatencyReset          types.Int64  `tfsdk:"rec_dns_latency_reset"`
	RecDnsLatencyListenOnSource types.String `tfsdk:"rec_dns_latency_listen_on_source"`
	RecDnsLatencyListenOnIp     types.String `tfsdk:"rec_dns_latency_listen_on_ip"`
	KpiMonitoredDomains         types.List   `tfsdk:"kpi_monitored_domains"`
}

var MemberTrafficCaptureRecDnsSettingAttrTypes = map[string]attr.Type{
	"rec_dns_latency_trigger_enable":   types.BoolType,
	"rec_dns_latency_threshold":        types.Int64Type,
	"rec_dns_latency_reset":            types.Int64Type,
	"rec_dns_latency_listen_on_source": types.StringType,
	"rec_dns_latency_listen_on_ip":     types.StringType,
	"kpi_monitored_domains":            types.ListType{ElemType: types.ObjectType{AttrTypes: MembertrafficcapturerecdnssettingKpiMonitoredDomainsAttrTypes}},
}

var MemberTrafficCaptureRecDnsSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"rec_dns_latency_trigger_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable triggering automated traffic capture based on recursive DNS latency.",
	},
	"rec_dns_latency_threshold": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Recursive DNS latency below which traffic capture will be triggered.",
	},
	"rec_dns_latency_reset": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Recursive DNS latency above which traffic capture will be stopped.",
	},
	"rec_dns_latency_listen_on_source": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("VIP_V4"),
		Validators: []validator.String{
			stringvalidator.OneOf("IP", "LAN2_V4", "LAN2_V6", "MGMT_V4", "MGMT_V6", "VIP_V4", "VIP_V6"),
		},
		MarkdownDescription: "The local IP DNS service is listen on ( for recursive DNS latency trigger).",
	},
	"rec_dns_latency_listen_on_ip": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The DNS listen-on IP address used if rec_dns_latency_listen_on_source is IP.",
	},
	"kpi_monitored_domains": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MembertrafficcapturerecdnssettingKpiMonitoredDomainsResourceSchemaAttributes,
		},
		Computed: true,
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "List of domains monitored by 'Recursive DNS Latency Threshold' trigger.",
	},
}

func ExpandMemberTrafficCaptureRecDnsSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberTrafficCaptureRecDnsSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberTrafficCaptureRecDnsSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberTrafficCaptureRecDnsSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberTrafficCaptureRecDnsSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberTrafficCaptureRecDnsSetting{
		RecDnsLatencyTriggerEnable:  flex.ExpandBoolPointer(m.RecDnsLatencyTriggerEnable),
		RecDnsLatencyThreshold:      flex.ExpandInt64Pointer(m.RecDnsLatencyThreshold),
		RecDnsLatencyReset:          flex.ExpandInt64Pointer(m.RecDnsLatencyReset),
		RecDnsLatencyListenOnSource: flex.ExpandStringPointer(m.RecDnsLatencyListenOnSource),
		RecDnsLatencyListenOnIp:     flex.ExpandStringPointerEmptyAsNil(m.RecDnsLatencyListenOnIp),
		KpiMonitoredDomains:         flex.ExpandFrameworkListNestedBlock(ctx, m.KpiMonitoredDomains, diags, ExpandMembertrafficcapturerecdnssettingKpiMonitoredDomains),
	}
	return to
}

func FlattenMemberTrafficCaptureRecDnsSetting(ctx context.Context, from *grid.MemberTrafficCaptureRecDnsSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberTrafficCaptureRecDnsSettingAttrTypes)
	}
	m := MemberTrafficCaptureRecDnsSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberTrafficCaptureRecDnsSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberTrafficCaptureRecDnsSettingModel) Flatten(ctx context.Context, from *grid.MemberTrafficCaptureRecDnsSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberTrafficCaptureRecDnsSettingModel{}
	}
	m.RecDnsLatencyTriggerEnable = types.BoolPointerValue(from.RecDnsLatencyTriggerEnable)
	m.RecDnsLatencyThreshold = flex.FlattenInt64Pointer(from.RecDnsLatencyThreshold)
	m.RecDnsLatencyReset = flex.FlattenInt64Pointer(from.RecDnsLatencyReset)
	m.RecDnsLatencyListenOnSource = flex.FlattenStringPointer(from.RecDnsLatencyListenOnSource)
	m.RecDnsLatencyListenOnIp = flex.FlattenStringPointer(from.RecDnsLatencyListenOnIp)
	m.KpiMonitoredDomains = flex.FlattenFrameworkListNestedBlock(ctx, from.KpiMonitoredDomains, MembertrafficcapturerecdnssettingKpiMonitoredDomainsAttrTypes, diags, FlattenMembertrafficcapturerecdnssettingKpiMonitoredDomains)
}

func (m *MemberTrafficCaptureRecDnsSettingModel) PutExpand(to *grid.MemberTrafficCaptureRecDnsSetting) *grid.MemberTrafficCaptureRecDnsSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberTrafficCaptureRecDnsSettingResourceSchemaAttributes {
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
