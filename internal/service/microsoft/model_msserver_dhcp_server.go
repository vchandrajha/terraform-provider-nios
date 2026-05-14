package microsoft

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/microsoft"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MsserverDhcpServerModel struct {
	UseLogin                   types.Bool   `tfsdk:"use_login"`
	LoginName                  types.String `tfsdk:"login_name"`
	LoginPassword              types.String `tfsdk:"login_password"`
	Managed                    types.Bool   `tfsdk:"managed"`
	NextSyncControl            types.String `tfsdk:"next_sync_control"`
	Status                     types.String `tfsdk:"status"`
	StatusLastUpdated          types.Int64  `tfsdk:"status_last_updated"`
	UseEnableMonitoring        types.Bool   `tfsdk:"use_enable_monitoring"`
	EnableMonitoring           types.Bool   `tfsdk:"enable_monitoring"`
	UseEnableInvalidMac        types.Bool   `tfsdk:"use_enable_invalid_mac"`
	EnableInvalidMac           types.Bool   `tfsdk:"enable_invalid_mac"`
	SupportsFailover           types.Bool   `tfsdk:"supports_failover"`
	UseSynchronizationMinDelay types.Bool   `tfsdk:"use_synchronization_min_delay"`
	SynchronizationMinDelay    types.Int64  `tfsdk:"synchronization_min_delay"`
}

var MsserverDhcpServerAttrTypes = map[string]attr.Type{
	"use_login":                     types.BoolType,
	"login_name":                    types.StringType,
	"login_password":                types.StringType,
	"managed":                       types.BoolType,
	"next_sync_control":             types.StringType,
	"status":                        types.StringType,
	"status_last_updated":           types.Int64Type,
	"use_enable_monitoring":         types.BoolType,
	"enable_monitoring":             types.BoolType,
	"use_enable_invalid_mac":        types.BoolType,
	"enable_invalid_mac":            types.BoolType,
	"supports_failover":             types.BoolType,
	"use_synchronization_min_delay": types.BoolType,
	"synchronization_min_delay":     types.Int64Type,
}

var MsserverDhcpServerResourceSchemaAttributes = map[string]schema.Attribute{
	"use_login": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag to override login name and password from the MS Server",
	},
	"login_name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Microsoft Server login name",
	},
	"login_password": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Sensitive:           true,
		MarkdownDescription: "Microsoft Server login password",
	},
	"managed": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "flag indicating if the DNS service is managed",
	},
	"next_sync_control": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("NONE", "START", "STOP"),
		},
		MarkdownDescription: "Defines what control to apply on the DNS server",
	},
	"status": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Status of the Microsoft DNS Service",
	},
	"status_last_updated": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Timestamp of the last update",
	},
	"use_enable_monitoring": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Override enable monitoring inherited from grid level",
	},
	"enable_monitoring": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag indicating if the DNS service is monitored and controlled",
	},
	"use_enable_invalid_mac": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Override setting for Enable Invalid Mac Address",
	},
	"enable_invalid_mac": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Enable Invalid Mac Address",
	},
	"supports_failover": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag indicating if the DHCP supports Failover",
	},
	"use_synchronization_min_delay": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag to override synchronization interval from the MS Server",
	},
	"synchronization_min_delay": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Minimum number of minutes between two synchronizations",
	},
}

func ExpandMsserverDhcpServer(ctx context.Context, o types.Object, diags *diag.Diagnostics) *microsoft.MsserverDhcpServer {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MsserverDhcpServerModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MsserverDhcpServerModel) Expand(ctx context.Context, diags *diag.Diagnostics) *microsoft.MsserverDhcpServer {
	if m == nil {
		return nil
	}
	to := &microsoft.MsserverDhcpServer{
		UseLogin:                   flex.ExpandBoolPointer(m.UseLogin),
		LoginName:                  flex.ExpandStringPointer(m.LoginName),
		LoginPassword:              flex.ExpandStringPointer(m.LoginPassword),
		Managed:                    flex.ExpandBoolPointer(m.Managed),
		NextSyncControl:            flex.ExpandStringPointer(m.NextSyncControl),
		UseEnableMonitoring:        flex.ExpandBoolPointer(m.UseEnableMonitoring),
		EnableMonitoring:           flex.ExpandBoolPointer(m.EnableMonitoring),
		UseEnableInvalidMac:        flex.ExpandBoolPointer(m.UseEnableInvalidMac),
		EnableInvalidMac:           flex.ExpandBoolPointer(m.EnableInvalidMac),
		UseSynchronizationMinDelay: flex.ExpandBoolPointer(m.UseSynchronizationMinDelay),
		SynchronizationMinDelay:    flex.ExpandInt64Pointer(m.SynchronizationMinDelay),
	}
	return to
}

func FlattenMsserverDhcpServer(ctx context.Context, from *microsoft.MsserverDhcpServer, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MsserverDhcpServerAttrTypes)
	}
	m := MsserverDhcpServerModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MsserverDhcpServerAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MsserverDhcpServerModel) Flatten(ctx context.Context, from *microsoft.MsserverDhcpServer, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MsserverDhcpServerModel{}
	}
	m.UseLogin = types.BoolPointerValue(from.UseLogin)
	m.LoginName = flex.FlattenStringPointer(from.LoginName)
	m.LoginPassword = flex.FlattenStringPointer(from.LoginPassword)
	m.Managed = types.BoolPointerValue(from.Managed)
	m.NextSyncControl = flex.FlattenStringPointer(from.NextSyncControl)
	m.Status = flex.FlattenStringPointer(from.Status)
	m.StatusLastUpdated = flex.FlattenInt64Pointer(from.StatusLastUpdated)
	m.UseEnableMonitoring = types.BoolPointerValue(from.UseEnableMonitoring)
	m.EnableMonitoring = types.BoolPointerValue(from.EnableMonitoring)
	m.UseEnableInvalidMac = types.BoolPointerValue(from.UseEnableInvalidMac)
	m.EnableInvalidMac = types.BoolPointerValue(from.EnableInvalidMac)
	m.SupportsFailover = types.BoolPointerValue(from.SupportsFailover)
	m.UseSynchronizationMinDelay = types.BoolPointerValue(from.UseSynchronizationMinDelay)
	m.SynchronizationMinDelay = flex.FlattenInt64Pointer(from.SynchronizationMinDelay)
}

func (m *MsserverDhcpServerModel) PutExpand(to *microsoft.MsserverDhcpServer) *microsoft.MsserverDhcpServer {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MsserverDhcpServerResourceSchemaAttributes {
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
