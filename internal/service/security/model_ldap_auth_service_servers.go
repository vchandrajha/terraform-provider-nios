package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type LdapAuthServiceServersModel struct {
	Address            types.String `tfsdk:"address"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	BaseDn             types.String `tfsdk:"base_dn"`
	BindPassword       types.String `tfsdk:"bind_password"`
	BindUserDn         types.String `tfsdk:"bind_user_dn"`
	Comment            types.String `tfsdk:"comment"`
	Disable            types.Bool   `tfsdk:"disable"`
	Encryption         types.String `tfsdk:"encryption"`
	Port               types.Int64  `tfsdk:"port"`
	UseMgmtPort        types.Bool   `tfsdk:"use_mgmt_port"`
	Version            types.String `tfsdk:"version"`
}

var LdapAuthServiceServersAttrTypes = map[string]attr.Type{
	"address":             types.StringType,
	"authentication_type": types.StringType,
	"base_dn":             types.StringType,
	"bind_password":       types.StringType,
	"bind_user_dn":        types.StringType,
	"comment":             types.StringType,
	"disable":             types.BoolType,
	"encryption":          types.StringType,
	"port":                types.Int64Type,
	"use_mgmt_port":       types.BoolType,
	"version":             types.StringType,
}

var LdapAuthServiceServersResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The IP address or FQDN of the LDAP server.",
	},
	"authentication_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ANONYMOUS", "AUTHENTICATED"),
		},
		Default:             stringdefault.StaticString("ANONYMOUS"),
		MarkdownDescription: "The authentication type for the LDAP server.",
	},
	"base_dn": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The base DN for the LDAP server.",
	},
	"bind_password": schema.StringAttribute{
		Optional:            true,
		Sensitive:           true,
		MarkdownDescription: "The user password for authentication.",
	},
	"bind_user_dn": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The user DN for authentication.",
	},
	"comment": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The LDAP descriptive comment.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the LDAP server is disabled.",
	},
	"encryption": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("NONE", "SSL"),
		},
		Default:             stringdefault.StaticString("SSL"),
		MarkdownDescription: "The LDAP server encryption type.",
	},
	"port": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(636),
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
		MarkdownDescription: "The LDAP server port.",
	},
	"use_mgmt_port": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the connection via the MGMT interface is allowed.",
	},
	"version": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("V2", "V3"),
		},
		Default:             stringdefault.StaticString("V3"),
		MarkdownDescription: "The LDAP server version.",
	},
}

func ExpandLdapAuthServiceServers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.LdapAuthServiceServers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m LdapAuthServiceServersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *LdapAuthServiceServersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.LdapAuthServiceServers {
	if m == nil {
		return nil
	}
	to := &security.LdapAuthServiceServers{
		Address:            flex.ExpandStringPointer(m.Address),
		AuthenticationType: flex.ExpandStringPointer(m.AuthenticationType),
		BaseDn:             flex.ExpandStringPointer(m.BaseDn),
		BindPassword:       flex.ExpandStringPointer(m.BindPassword),
		BindUserDn:         flex.ExpandStringPointer(m.BindUserDn),
		Comment:            flex.ExpandStringPointer(m.Comment),
		Disable:            flex.ExpandBoolPointer(m.Disable),
		Encryption:         flex.ExpandStringPointer(m.Encryption),
		Port:               flex.ExpandInt64Pointer(m.Port),
		UseMgmtPort:        flex.ExpandBoolPointer(m.UseMgmtPort),
		Version:            flex.ExpandStringPointer(m.Version),
	}
	return to
}

func FlattenLdapAuthServiceServers(ctx context.Context, from *security.LdapAuthServiceServers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(LdapAuthServiceServersAttrTypes)
	}
	m := LdapAuthServiceServersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, LdapAuthServiceServersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *LdapAuthServiceServersModel) Flatten(ctx context.Context, from *security.LdapAuthServiceServers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = LdapAuthServiceServersModel{}
	}
	m.Address = flex.FlattenStringPointer(from.Address)
	m.AuthenticationType = flex.FlattenStringPointer(from.AuthenticationType)
	m.BaseDn = flex.FlattenStringPointer(from.BaseDn)
	m.BindUserDn = flex.FlattenStringPointer(from.BindUserDn)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.Encryption = flex.FlattenStringPointer(from.Encryption)
	m.Port = flex.FlattenInt64Pointer(from.Port)
	m.UseMgmtPort = types.BoolPointerValue(from.UseMgmtPort)
	m.Version = flex.FlattenStringPointer(from.Version)
}

func (m *LdapAuthServiceServersModel) PutExpand(to *security.LdapAuthServiceServers) *security.LdapAuthServiceServers {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range LdapAuthServiceServersResourceSchemaAttributes {
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
