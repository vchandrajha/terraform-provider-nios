package misc

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/misc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type SyslogEndpointSyslogServersModel struct {
	Address             iptypes.IPAddress `tfsdk:"address"`
	ConnectionType      types.String      `tfsdk:"connection_type"`
	Port                types.Int64       `tfsdk:"port"`
	Hostname            types.String      `tfsdk:"hostname"`
	Format              types.String      `tfsdk:"format"`
	Facility            types.String      `tfsdk:"facility"`
	Severity            types.String      `tfsdk:"severity"`
	Certificate         types.Object      `tfsdk:"certificate"`
	CertificateToken    types.String      `tfsdk:"certificate_token"`
	CertificateFilePath types.String      `tfsdk:"certificate_file_path"`
}

var SyslogEndpointSyslogServersAttrTypes = map[string]attr.Type{
	"address":               iptypes.IPAddressType{},
	"connection_type":       types.StringType,
	"port":                  types.Int64Type,
	"hostname":              types.StringType,
	"format":                types.StringType,
	"facility":              types.StringType,
	"severity":              types.StringType,
	"certificate":           types.ObjectType{AttrTypes: SyslogEndpointSyslogServersCertificateAttrTypes},
	"certificate_token":     types.StringType,
	"certificate_file_path": types.StringType,
}

var SyslogEndpointSyslogServersResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Required:            true,
		CustomType:          iptypes.IPAddressType{},
		MarkdownDescription: "Syslog Server IP address",
	},
	"connection_type": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("stcp", "udp", "tls"),
		},
		Default:             stringdefault.StaticString("udp"),
		MarkdownDescription: "Connection type values",
	},
	"port": schema.Int64Attribute{
		Computed:            true,
		Optional:            true,
		Default:             int64default.StaticInt64(514),
		MarkdownDescription: "The port this server listens on.",
	},
	"hostname": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("HOSTNAME", "FQDN", "IP_ADDRESS"),
		},
		Default:             stringdefault.StaticString("HOSTNAME"),
		MarkdownDescription: "List of hostnames",
	},
	"format": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("formatted", "raw"),
		},
		Default:             stringdefault.StaticString("raw"),
		MarkdownDescription: "Format values for syslog endpoint server",
	},
	"facility": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7"),
		},
		Default:             stringdefault.StaticString("local0"),
		MarkdownDescription: "Facility values for syslog endpoint server",
	},
	"severity": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("alert", "critic", "debug", "emerg", "err", "info", "notice", "warning"),
		},
		Default:             stringdefault.StaticString("debug"),
		MarkdownDescription: "Severity values for syslog endpoint server.",
	},
	"certificate": schema.SingleNestedAttribute{
		Computed:            true,
		Attributes:          SyslogEndpointSyslogServersCertificateResourceSchemaAttributes,
		MarkdownDescription: "Reference for creating syslog endpoint server.",
	},
	"certificate_token": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The token returned by the uploadinit function call in object fileop.",
	},
	"certificate_file_path": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The file path to the certificate.",
	},
}

func ExpandSyslogEndpointSyslogServers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *misc.SyslogEndpointSyslogServers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m SyslogEndpointSyslogServersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *SyslogEndpointSyslogServersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *misc.SyslogEndpointSyslogServers {
	if m == nil {
		return nil
	}
	to := &misc.SyslogEndpointSyslogServers{
		Address:          flex.ExpandIPAddress(m.Address),
		ConnectionType:   flex.ExpandStringPointer(m.ConnectionType),
		Port:             flex.ExpandInt64Pointer(m.Port),
		Hostname:         flex.ExpandStringPointer(m.Hostname),
		Format:           flex.ExpandStringPointer(m.Format),
		Facility:         flex.ExpandStringPointer(m.Facility),
		Severity:         flex.ExpandStringPointer(m.Severity),
		CertificateToken: flex.ExpandStringPointer(m.CertificateToken),
	}
	return to
}

func FlattenSyslogEndpointSyslogServers(ctx context.Context, from *misc.SyslogEndpointSyslogServers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SyslogEndpointSyslogServersAttrTypes)
	}
	m := SyslogEndpointSyslogServersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, SyslogEndpointSyslogServersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SyslogEndpointSyslogServersModel) Flatten(ctx context.Context, from *misc.SyslogEndpointSyslogServers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SyslogEndpointSyslogServersModel{}
	}
	m.Address = flex.FlattenIPAddress(from.Address)
	m.ConnectionType = flex.FlattenStringPointer(from.ConnectionType)
	m.Port = flex.FlattenInt64Pointer(from.Port)
	m.Hostname = flex.FlattenStringPointer(from.Hostname)
	m.Format = flex.FlattenStringPointer(from.Format)
	m.Facility = flex.FlattenStringPointer(from.Facility)
	m.Severity = flex.FlattenStringPointer(from.Severity)
	m.Certificate = flattenCertificate(ctx, from.Certificate, diags)
	m.CertificateToken = flex.FlattenStringPointer(from.CertificateToken)
	m.CertificateFilePath = types.StringNull() // This field is write-only and should not be set from API data
}

// flattenCertificate handles both string and object types returned by the API.
// When the API returns an object (for stcp connections), it extracts the certificate fields.
func flattenCertificate(ctx context.Context, cert interface{}, diags *diag.Diagnostics) types.Object {
	if cert == nil {
		return types.ObjectNull(SyslogEndpointSyslogServersCertificateAttrTypes)
	}

	switch v := cert.(type) {
	case string:
		// If certificate is a string, put it in the ref field
		m := SyslogEndpointSyslogServersCertificateModel{
			Ref: types.StringValue(v),
		}
		obj, d := types.ObjectValueFrom(ctx, SyslogEndpointSyslogServersCertificateAttrTypes, m)
		diags.Append(d...)
		return obj
	case map[string]interface{}:
		// API returns certificate as an object with _ref, issuer, serial, etc.
		m := SyslogEndpointSyslogServersCertificateModel{}
		if ref, ok := v["_ref"].(string); ok {
			m.Ref = types.StringValue(ref)
		}
		if issuer, ok := v["issuer"].(string); ok {
			m.Issuer = types.StringValue(issuer)
		}
		if serial, ok := v["serial"].(string); ok {
			m.Serial = types.StringValue(serial)
		}
		if subject, ok := v["subject"].(string); ok {
			m.Subject = types.StringValue(subject)
		}
		// Handle valid_not_after - can be float64 from JSON unmarshaling
		if validNotAfter, ok := v["valid_not_after"].(float64); ok {
			m.ValidNotAfter = types.Int64Value(int64(validNotAfter))
		} else if validNotAfter, ok := v["valid_not_after"].(int64); ok {
			m.ValidNotAfter = types.Int64Value(validNotAfter)
		}
		// Handle valid_not_before - can be float64 from JSON unmarshaling
		if validNotBefore, ok := v["valid_not_before"].(float64); ok {
			m.ValidNotBefore = types.Int64Value(int64(validNotBefore))
		} else if validNotBefore, ok := v["valid_not_before"].(int64); ok {
			m.ValidNotBefore = types.Int64Value(validNotBefore)
		}
		obj, d := types.ObjectValueFrom(ctx, SyslogEndpointSyslogServersCertificateAttrTypes, m)
		diags.Append(d...)
		return obj
	default:
		return types.ObjectNull(SyslogEndpointSyslogServersCertificateAttrTypes)
	}
}

func (m *SyslogEndpointSyslogServersModel) PutExpand(to *misc.SyslogEndpointSyslogServers) *misc.SyslogEndpointSyslogServers {
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

	for field, attr := range SyslogEndpointSyslogServersResourceSchemaAttributes {
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
