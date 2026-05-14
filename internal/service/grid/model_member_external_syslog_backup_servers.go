package grid

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

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type MemberExternalSyslogBackupServersModel struct {
	AddressOrFqdn types.String `tfsdk:"address_or_fqdn"`
	DirectoryPath types.String `tfsdk:"directory_path"`
	Enable        types.Bool   `tfsdk:"enable"`
	Password      types.String `tfsdk:"password"`
	Port          types.Int64  `tfsdk:"port"`
	Protocol      types.String `tfsdk:"protocol"`
	Username      types.String `tfsdk:"username"`
}

var MemberExternalSyslogBackupServersAttrTypes = map[string]attr.Type{
	"address_or_fqdn": types.StringType,
	"directory_path":  types.StringType,
	"enable":          types.BoolType,
	"password":        types.StringType,
	"port":            types.Int64Type,
	"protocol":        types.StringType,
	"username":        types.StringType,
}

var MemberExternalSyslogBackupServersResourceSchemaAttributes = map[string]schema.Attribute{
	"address_or_fqdn": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The IPv4 or IPv6 address or FQDN of the backup syslog server.",
	},
	"directory_path": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The directory path for the replication of the rotated syslog files.",
	},
	"enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
		MarkdownDescription: "If set to True, the syslog backup server is enabled.",
	},
	"password": schema.StringAttribute{
		Required:            true,
		Sensitive:           true,
		MarkdownDescription: "The password of the backup syslog server.",
	},
	"port": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(22),
		Validators: []validator.Int64{
			int64validator.Between(0, 65535),
		},
		MarkdownDescription: "The port used to connect to the backup syslog server.",
	},
	"protocol": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("SCP"),
		Validators: []validator.String{
			stringvalidator.OneOf("FTP", "SCP"),
		},
		MarkdownDescription: "The transport protocol used to connect to the backup syslog server.",
	},
	"username": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The username of the backup syslog server.",
	},
}

func ExpandMemberExternalSyslogBackupServers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberExternalSyslogBackupServers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberExternalSyslogBackupServersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberExternalSyslogBackupServersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberExternalSyslogBackupServers {
	if m == nil {
		return nil
	}
	to := &grid.MemberExternalSyslogBackupServers{
		AddressOrFqdn: flex.ExpandStringPointer(m.AddressOrFqdn),
		DirectoryPath: flex.ExpandStringPointer(m.DirectoryPath),
		Enable:        flex.ExpandBoolPointer(m.Enable),
		Password:      flex.ExpandStringPointer(m.Password),
		Port:          flex.ExpandInt64Pointer(m.Port),
		Protocol:      flex.ExpandStringPointer(m.Protocol),
		Username:      flex.ExpandStringPointer(m.Username),
	}
	return to
}

func FlattenMemberExternalSyslogBackupServers(ctx context.Context, from *grid.MemberExternalSyslogBackupServers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberExternalSyslogBackupServersAttrTypes)
	}
	m := MemberExternalSyslogBackupServersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberExternalSyslogBackupServersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberExternalSyslogBackupServersModel) Flatten(ctx context.Context, from *grid.MemberExternalSyslogBackupServers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberExternalSyslogBackupServersModel{}
	}
	m.AddressOrFqdn = flex.FlattenStringPointer(from.AddressOrFqdn)
	m.DirectoryPath = flex.FlattenStringPointer(from.DirectoryPath)
	m.Enable = types.BoolPointerValue(from.Enable)
	m.Port = flex.FlattenInt64Pointer(from.Port)
	m.Protocol = flex.FlattenStringPointer(from.Protocol)
	m.Username = flex.FlattenStringPointer(from.Username)
	m.Password = flex.FlattenStringPointer(from.Password)
}

func (m *MemberExternalSyslogBackupServersModel) PutExpand(to *grid.MemberExternalSyslogBackupServers) *grid.MemberExternalSyslogBackupServers {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberExternalSyslogBackupServersResourceSchemaAttributes {
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
