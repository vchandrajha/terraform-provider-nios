package grid

import (
	"context"
	"reflect"
	"strings"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberTrafficCaptureAuthDnsSettingModel struct {
	AuthDnsLatencyTriggerEnable  types.Bool   `tfsdk:"auth_dns_latency_trigger_enable"`
	AuthDnsLatencyThreshold      types.Int64  `tfsdk:"auth_dns_latency_threshold"`
	AuthDnsLatencyReset          types.Int64  `tfsdk:"auth_dns_latency_reset"`
	AuthDnsLatencyListenOnSource types.String `tfsdk:"auth_dns_latency_listen_on_source"`
	AuthDnsLatencyListenOnIp     types.String `tfsdk:"auth_dns_latency_listen_on_ip"`
}

var MemberTrafficCaptureAuthDnsSettingAttrTypes = map[string]attr.Type{
	"auth_dns_latency_trigger_enable":   types.BoolType,
	"auth_dns_latency_threshold":        types.Int64Type,
	"auth_dns_latency_reset":            types.Int64Type,
	"auth_dns_latency_listen_on_source": types.StringType,
	"auth_dns_latency_listen_on_ip":     types.StringType,
}

var MemberTrafficCaptureAuthDnsSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"auth_dns_latency_trigger_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enabling trigger automated traffic capture based on authoritative DNS latency.",
	},
	"auth_dns_latency_threshold": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Authoritative DNS latency below which traffic capture will be triggered.",
	},
	"auth_dns_latency_reset": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Authoritative DNS latency above which traffic capture will be stopped.",
	},
	"auth_dns_latency_listen_on_source": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("VIP_V4"),
		Validators: []validator.String{
			stringvalidator.OneOf("IP", "LAN2_V4", "LAN2_V6", "MGMT_V4", "MGMT_V6", "VIP_V4", "VIP_V6"),
		},
		MarkdownDescription: "The local IP DNS service is listen on (for authoritative DNS latency trigger).",
	},
	"auth_dns_latency_listen_on_ip": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The DNS listen-on IP address used if auth_dns_latency_on_source is IP.",
	},
}

func ExpandMemberTrafficCaptureAuthDnsSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberTrafficCaptureAuthDnsSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberTrafficCaptureAuthDnsSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberTrafficCaptureAuthDnsSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberTrafficCaptureAuthDnsSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberTrafficCaptureAuthDnsSetting{
		AuthDnsLatencyTriggerEnable:  flex.ExpandBoolPointer(m.AuthDnsLatencyTriggerEnable),
		AuthDnsLatencyThreshold:      flex.ExpandInt64Pointer(m.AuthDnsLatencyThreshold),
		AuthDnsLatencyReset:          flex.ExpandInt64Pointer(m.AuthDnsLatencyReset),
		AuthDnsLatencyListenOnSource: flex.ExpandStringPointer(m.AuthDnsLatencyListenOnSource),
		AuthDnsLatencyListenOnIp:     flex.ExpandStringPointerEmptyAsNil(m.AuthDnsLatencyListenOnIp),
	}
	return to
}

func FlattenMemberTrafficCaptureAuthDnsSetting(ctx context.Context, from *grid.MemberTrafficCaptureAuthDnsSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberTrafficCaptureAuthDnsSettingAttrTypes)
	}
	m := MemberTrafficCaptureAuthDnsSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberTrafficCaptureAuthDnsSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberTrafficCaptureAuthDnsSettingModel) Flatten(ctx context.Context, from *grid.MemberTrafficCaptureAuthDnsSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberTrafficCaptureAuthDnsSettingModel{}
	}
	m.AuthDnsLatencyTriggerEnable = types.BoolPointerValue(from.AuthDnsLatencyTriggerEnable)
	m.AuthDnsLatencyThreshold = flex.FlattenInt64Pointer(from.AuthDnsLatencyThreshold)
	m.AuthDnsLatencyReset = flex.FlattenInt64Pointer(from.AuthDnsLatencyReset)
	m.AuthDnsLatencyListenOnSource = flex.FlattenStringPointer(from.AuthDnsLatencyListenOnSource)
	m.AuthDnsLatencyListenOnIp = flex.FlattenStringPointer(from.AuthDnsLatencyListenOnIp)
}

func (m *MemberTrafficCaptureAuthDnsSettingModel) PutExpand(to *grid.MemberTrafficCaptureAuthDnsSetting) *grid.MemberTrafficCaptureAuthDnsSetting {
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

	for field, attr := range MemberTrafficCaptureAuthDnsSettingResourceSchemaAttributes {
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
