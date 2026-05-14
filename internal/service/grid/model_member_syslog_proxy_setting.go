package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

type MemberSyslogProxySettingModel struct {
	Enable     types.Bool  `tfsdk:"enable"`
	TcpEnable  types.Bool  `tfsdk:"tcp_enable"`
	TcpPort    types.Int64 `tfsdk:"tcp_port"`
	UdpEnable  types.Bool  `tfsdk:"udp_enable"`
	UdpPort    types.Int64 `tfsdk:"udp_port"`
	ClientAcls types.List  `tfsdk:"client_acls"`
}

var MemberSyslogProxySettingAttrTypes = map[string]attr.Type{
	"enable":      types.BoolType,
	"tcp_enable":  types.BoolType,
	"tcp_port":    types.Int64Type,
	"udp_enable":  types.BoolType,
	"udp_port":    types.Int64Type,
	"client_acls": types.ListType{ElemType: types.ObjectType{AttrTypes: MembersyslogproxysettingClientAclsAttrTypes}},
}

// Defaults removed due to incorrect default values
var MemberSyslogProxySettingResourceSchemaAttributes = map[string]schema.Attribute{
	"enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "If set to True, the member receives syslog messages from specified devices, such as syslog servers and routers, and then forwards these messages to an external syslog server.",
	},
	"tcp_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If set to True, the appliance can receive messages from other devices via TCP.",
	},
	"tcp_port": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The TCP port the appliance must listen on.",
	},
	"udp_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If set to True, the appliance can receive messages from other devices via UDP.",
	},
	"udp_port": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		//Default:             int64default.StaticInt64(514),
		MarkdownDescription: "The UDP port the appliance must listen on.",
	},
	"client_acls": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MembersyslogproxysettingClientAclsResourceSchemaAttributes,
		},
		Computed: true,
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "This list controls the IP addresses and networks that are allowed to access the syslog proxy.",
	},
}

func ExpandMemberSyslogProxySetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberSyslogProxySetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberSyslogProxySettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberSyslogProxySettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberSyslogProxySetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberSyslogProxySetting{
		Enable:     flex.ExpandBoolPointer(m.Enable),
		TcpEnable:  flex.ExpandBoolPointer(m.TcpEnable),
		TcpPort:    flex.ExpandInt64Pointer(m.TcpPort),
		UdpEnable:  flex.ExpandBoolPointer(m.UdpEnable),
		UdpPort:    flex.ExpandInt64Pointer(m.UdpPort),
		ClientAcls: flex.ExpandFrameworkListNestedBlock(ctx, m.ClientAcls, diags, ExpandMembersyslogproxysettingClientAcls),
	}
	return to
}

func FlattenMemberSyslogProxySetting(ctx context.Context, from *grid.MemberSyslogProxySetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberSyslogProxySettingAttrTypes)
	}
	m := MemberSyslogProxySettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberSyslogProxySettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberSyslogProxySettingModel) Flatten(ctx context.Context, from *grid.MemberSyslogProxySetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberSyslogProxySettingModel{}
	}
	m.Enable = types.BoolPointerValue(from.Enable)
	m.TcpEnable = types.BoolPointerValue(from.TcpEnable)
	m.TcpPort = flex.FlattenInt64Pointer(from.TcpPort)
	m.UdpEnable = types.BoolPointerValue(from.UdpEnable)
	m.UdpPort = flex.FlattenInt64Pointer(from.UdpPort)
	m.ClientAcls = flex.FlattenFrameworkListNestedBlock(ctx, from.ClientAcls, MembersyslogproxysettingClientAclsAttrTypes, diags, FlattenMembersyslogproxysettingClientAcls)
}

func (m *MemberSyslogProxySettingModel) PutExpand(to *grid.MemberSyslogProxySetting) *grid.MemberSyslogProxySetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberSyslogProxySettingResourceSchemaAttributes {
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
