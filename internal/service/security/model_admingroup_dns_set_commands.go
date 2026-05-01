package security

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupDnsSetCommandsModel struct {
	SetDns                          types.Bool `tfsdk:"set_dns"`
	SetDnsRrl                       types.Bool `tfsdk:"set_dns_rrl"`
	SetEnableDnstap                 types.Bool `tfsdk:"set_enable_dnstap"`
	SetEnableMatchRecursiveOnly     types.Bool `tfsdk:"set_enable_match_recursive_only"`
	SetExtraDnsNameValidations      types.Bool `tfsdk:"set_extra_dns_name_validations"`
	SetLogGuestLookups              types.Bool `tfsdk:"set_log_guest_lookups"`
	SetMaxRecursionDepth            types.Bool `tfsdk:"set_max_recursion_depth"`
	SetMaxRecursionQueries          types.Bool `tfsdk:"set_max_recursion_queries"`
	SetMonitor                      types.Bool `tfsdk:"set_monitor"`
	SetMsDnsReportsSyncInterval     types.Bool `tfsdk:"set_ms_dns_reports_sync_interval"`
	SetMsStickyIp                   types.Bool `tfsdk:"set_ms_sticky_ip"`
	SetRestartAnycastWithDnsRestart types.Bool `tfsdk:"set_restart_anycast_with_dns_restart"`
	SetRpzAddSoa                    types.Bool `tfsdk:"set_rpz_add_soa"`
	SetDnsAccel                     types.Bool `tfsdk:"set_dns_accel"`
	SetDnsAccelDebug                types.Bool `tfsdk:"set_dns_accel_debug"`
	SetDnsAutoGen                   types.Bool `tfsdk:"set_dns_auto_gen"`
	SetAllowQueryDomain             types.Bool `tfsdk:"set_allow_query_domain"`
	EnableAll                       types.Bool `tfsdk:"enable_all"`
	DisableAll                      types.Bool `tfsdk:"disable_all"`
}

var AdmingroupDnsSetCommandsAttrTypes = map[string]attr.Type{
	"set_dns":                              types.BoolType,
	"set_dns_rrl":                          types.BoolType,
	"set_enable_dnstap":                    types.BoolType,
	"set_enable_match_recursive_only":      types.BoolType,
	"set_extra_dns_name_validations":       types.BoolType,
	"set_log_guest_lookups":                types.BoolType,
	"set_max_recursion_depth":              types.BoolType,
	"set_max_recursion_queries":            types.BoolType,
	"set_monitor":                          types.BoolType,
	"set_ms_dns_reports_sync_interval":     types.BoolType,
	"set_ms_sticky_ip":                     types.BoolType,
	"set_restart_anycast_with_dns_restart": types.BoolType,
	"set_rpz_add_soa":                      types.BoolType,
	"set_dns_accel":                        types.BoolType,
	"set_dns_accel_debug":                  types.BoolType,
	"set_dns_auto_gen":                     types.BoolType,
	"set_allow_query_domain":               types.BoolType,
	"enable_all":                           types.BoolType,
	"disable_all":                          types.BoolType,
}

var AdmingroupDnsSetCommandsResourceSchemaAttributes = map[string]schema.Attribute{
	"set_dns": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_dns_rrl": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_enable_dnstap": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_enable_match_recursive_only": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_extra_dns_name_validations": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_log_guest_lookups": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_max_recursion_depth": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_max_recursion_queries": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_monitor": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_ms_dns_reports_sync_interval": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_ms_sticky_ip": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_restart_anycast_with_dns_restart": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_rpz_add_soa": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_dns_accel": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_dns_accel_debug": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_dns_auto_gen": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_allow_query_domain": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"enable_all": schema.BoolAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then enable all fields",
	},
	"disable_all": schema.BoolAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then disable all fields",
	},
}

func ExpandAdmingroupDnsSetCommands(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupDnsSetCommands {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupDnsSetCommandsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupDnsSetCommandsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupDnsSetCommands {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupDnsSetCommands{
		SetDns:                          flex.ExpandBoolPointer(m.SetDns),
		SetDnsRrl:                       flex.ExpandBoolPointer(m.SetDnsRrl),
		SetEnableDnstap:                 flex.ExpandBoolPointer(m.SetEnableDnstap),
		SetEnableMatchRecursiveOnly:     flex.ExpandBoolPointer(m.SetEnableMatchRecursiveOnly),
		SetExtraDnsNameValidations:      flex.ExpandBoolPointer(m.SetExtraDnsNameValidations),
		SetLogGuestLookups:              flex.ExpandBoolPointer(m.SetLogGuestLookups),
		SetMaxRecursionDepth:            flex.ExpandBoolPointer(m.SetMaxRecursionDepth),
		SetMaxRecursionQueries:          flex.ExpandBoolPointer(m.SetMaxRecursionQueries),
		SetMonitor:                      flex.ExpandBoolPointer(m.SetMonitor),
		SetMsDnsReportsSyncInterval:     flex.ExpandBoolPointer(m.SetMsDnsReportsSyncInterval),
		SetMsStickyIp:                   flex.ExpandBoolPointer(m.SetMsStickyIp),
		SetRestartAnycastWithDnsRestart: flex.ExpandBoolPointer(m.SetRestartAnycastWithDnsRestart),
		SetRpzAddSoa:                    flex.ExpandBoolPointer(m.SetRpzAddSoa),
		SetDnsAccel:                     flex.ExpandBoolPointer(m.SetDnsAccel),
		SetDnsAccelDebug:                flex.ExpandBoolPointer(m.SetDnsAccelDebug),
		SetDnsAutoGen:                   flex.ExpandBoolPointer(m.SetDnsAutoGen),
		SetAllowQueryDomain:             flex.ExpandBoolPointer(m.SetAllowQueryDomain),
	}
	return to
}

func FlattenAdmingroupDnsSetCommands(ctx context.Context, from *security.AdmingroupDnsSetCommands, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupDnsSetCommandsAttrTypes)
	}
	m := AdmingroupDnsSetCommandsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupDnsSetCommandsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupDnsSetCommandsModel) Flatten(ctx context.Context, from *security.AdmingroupDnsSetCommands, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupDnsSetCommandsModel{}
	}
	m.SetDns = types.BoolPointerValue(from.SetDns)
	m.SetDnsRrl = types.BoolPointerValue(from.SetDnsRrl)
	m.SetEnableDnstap = types.BoolPointerValue(from.SetEnableDnstap)
	m.SetEnableMatchRecursiveOnly = types.BoolPointerValue(from.SetEnableMatchRecursiveOnly)
	m.SetExtraDnsNameValidations = types.BoolPointerValue(from.SetExtraDnsNameValidations)
	m.SetLogGuestLookups = types.BoolPointerValue(from.SetLogGuestLookups)
	m.SetMaxRecursionDepth = types.BoolPointerValue(from.SetMaxRecursionDepth)
	m.SetMaxRecursionQueries = types.BoolPointerValue(from.SetMaxRecursionQueries)
	m.SetMonitor = types.BoolPointerValue(from.SetMonitor)
	m.SetMsDnsReportsSyncInterval = types.BoolPointerValue(from.SetMsDnsReportsSyncInterval)
	m.SetMsStickyIp = types.BoolPointerValue(from.SetMsStickyIp)
	m.SetRestartAnycastWithDnsRestart = types.BoolPointerValue(from.SetRestartAnycastWithDnsRestart)
	m.SetRpzAddSoa = types.BoolPointerValue(from.SetRpzAddSoa)
	m.SetDnsAccel = types.BoolPointerValue(from.SetDnsAccel)
	m.SetDnsAccelDebug = types.BoolPointerValue(from.SetDnsAccelDebug)
	m.SetDnsAutoGen = types.BoolPointerValue(from.SetDnsAutoGen)
	m.SetAllowQueryDomain = types.BoolPointerValue(from.SetAllowQueryDomain)
	m.EnableAll = types.BoolPointerValue(from.EnableAll)
	m.DisableAll = types.BoolPointerValue(from.DisableAll)
}

func (m *AdmingroupDnsSetCommandsModel) PutExpand(to *security.AdmingroupDnsSetCommands) *security.AdmingroupDnsSetCommands {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AdmingroupDnsSetCommandsResourceSchemaAttributes {
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
