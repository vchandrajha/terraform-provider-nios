package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupTroubleShootingToplevelCommandsModel struct {
	Ping           types.Bool `tfsdk:"ping"`
	Ping6          types.Bool `tfsdk:"ping6"`
	Strace         types.Bool `tfsdk:"strace"`
	Traceroute     types.Bool `tfsdk:"traceroute"`
	TrafficCapture types.Bool `tfsdk:"traffic_capture"`
	Dig            types.Bool `tfsdk:"dig"`
	Rotate         types.Bool `tfsdk:"rotate"`
	Snmpwalk       types.Bool `tfsdk:"snmpwalk"`
	Snmpget        types.Bool `tfsdk:"snmpget"`
	Console        types.Bool `tfsdk:"console"`
	Tracepath      types.Bool `tfsdk:"tracepath"`
	EnableAll      types.Bool `tfsdk:"enable_all"`
	DisableAll     types.Bool `tfsdk:"disable_all"`
}

var AdmingroupTroubleShootingToplevelCommandsAttrTypes = map[string]attr.Type{
	"ping":            types.BoolType,
	"ping6":           types.BoolType,
	"strace":          types.BoolType,
	"traceroute":      types.BoolType,
	"traffic_capture": types.BoolType,
	"dig":             types.BoolType,
	"rotate":          types.BoolType,
	"snmpwalk":        types.BoolType,
	"snmpget":         types.BoolType,
	"console":         types.BoolType,
	"tracepath":       types.BoolType,
	"enable_all":      types.BoolType,
	"disable_all":     types.BoolType,
}

var AdmingroupTroubleShootingToplevelCommandsResourceSchemaAttributes = map[string]schema.Attribute{
	"ping": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"ping6": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"strace": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"traceroute": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"traffic_capture": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"dig": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"rotate": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"snmpwalk": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"snmpget": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"console": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"tracepath": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"enable_all": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "If True then enable all fields",
	},
	"disable_all": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "If True then disable all fields",
	},
}

func ExpandAdmingroupTroubleShootingToplevelCommands(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupTroubleShootingToplevelCommands {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupTroubleShootingToplevelCommandsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupTroubleShootingToplevelCommandsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupTroubleShootingToplevelCommands {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupTroubleShootingToplevelCommands{
		Ping:           flex.ExpandBoolPointer(m.Ping),
		Ping6:          flex.ExpandBoolPointer(m.Ping6),
		Strace:         flex.ExpandBoolPointer(m.Strace),
		Traceroute:     flex.ExpandBoolPointer(m.Traceroute),
		TrafficCapture: flex.ExpandBoolPointer(m.TrafficCapture),
		Dig:            flex.ExpandBoolPointer(m.Dig),
		Rotate:         flex.ExpandBoolPointer(m.Rotate),
		Snmpwalk:       flex.ExpandBoolPointer(m.Snmpwalk),
		Snmpget:        flex.ExpandBoolPointer(m.Snmpget),
		Console:        flex.ExpandBoolPointer(m.Console),
		Tracepath:      flex.ExpandBoolPointer(m.Tracepath),
	}
	return to
}

func FlattenAdmingroupTroubleShootingToplevelCommands(ctx context.Context, from *security.AdmingroupTroubleShootingToplevelCommands, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupTroubleShootingToplevelCommandsAttrTypes)
	}
	m := AdmingroupTroubleShootingToplevelCommandsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupTroubleShootingToplevelCommandsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupTroubleShootingToplevelCommandsModel) Flatten(ctx context.Context, from *security.AdmingroupTroubleShootingToplevelCommands, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupTroubleShootingToplevelCommandsModel{}
	}
	m.Ping = types.BoolPointerValue(from.Ping)
	m.Ping6 = types.BoolPointerValue(from.Ping6)
	m.Strace = types.BoolPointerValue(from.Strace)
	m.Traceroute = types.BoolPointerValue(from.Traceroute)
	m.TrafficCapture = types.BoolPointerValue(from.TrafficCapture)
	m.Dig = types.BoolPointerValue(from.Dig)
	m.Rotate = types.BoolPointerValue(from.Rotate)
	m.Snmpwalk = types.BoolPointerValue(from.Snmpwalk)
	m.Snmpget = types.BoolPointerValue(from.Snmpget)
	m.Console = types.BoolPointerValue(from.Console)
	m.Tracepath = types.BoolPointerValue(from.Tracepath)
	m.EnableAll = types.BoolPointerValue(from.EnableAll)
	m.DisableAll = types.BoolPointerValue(from.DisableAll)
}

func (m *AdmingroupTroubleShootingToplevelCommandsModel) PutExpand(to *security.AdmingroupTroubleShootingToplevelCommands) *security.AdmingroupTroubleShootingToplevelCommands {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AdmingroupTroubleShootingToplevelCommandsResourceSchemaAttributes {
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
