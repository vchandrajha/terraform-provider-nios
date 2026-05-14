package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type NetworkviewRemoteReverseZonesModel struct {
	Fqdn                types.String      `tfsdk:"fqdn"`
	ServerAddress       iptypes.IPAddress `tfsdk:"server_address"`
	GssTsigDnsPrincipal types.String      `tfsdk:"gss_tsig_dns_principal"`
	GssTsigDomain       types.String      `tfsdk:"gss_tsig_domain"`
	TsigKey             types.String      `tfsdk:"tsig_key"`
	TsigKeyAlg          types.String      `tfsdk:"tsig_key_alg"`
	TsigKeyName         types.String      `tfsdk:"tsig_key_name"`
	KeyType             types.String      `tfsdk:"key_type"`
}

var NetworkviewRemoteReverseZonesAttrTypes = map[string]attr.Type{
	"fqdn":                   types.StringType,
	"server_address":         iptypes.IPAddressType{},
	"gss_tsig_dns_principal": types.StringType,
	"gss_tsig_domain":        types.StringType,
	"tsig_key":               types.StringType,
	"tsig_key_alg":           types.StringType,
	"tsig_key_name":          types.StringType,
	"key_type":               types.StringType,
}

var NetworkviewRemoteReverseZonesResourceSchemaAttributes = map[string]schema.Attribute{
	"fqdn": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.IsValidDomainName(),
		},
		MarkdownDescription: "The FQDN of the remote server.",
	},
	"server_address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Required:            true,
		MarkdownDescription: "The remote server IP address.",
	},
	"gss_tsig_dns_principal": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The principal name in which GSS-TSIG for dynamic updates is enabled.",
	},
	"gss_tsig_domain": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The domain in which GSS-TSIG for dynamic updates is enabled.",
	},
	"tsig_key": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The TSIG key value.",
	},
	"tsig_key_alg": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("HMAC-MD5", "HMAC-SHA256"),
		},
		MarkdownDescription: "The TSIG key alorithm name.",
	},
	"tsig_key_name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The name of the TSIG key. The key name entered here must match the TSIG key name on the external name server.",
	},
	"key_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("GSS-TSIG", "NONE", "TSIG"),
		},
		Default:             stringdefault.StaticString("NONE"),
		MarkdownDescription: "The key type to be used.",
	},
}

func ExpandNetworkviewRemoteReverseZones(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworkviewRemoteReverseZones {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworkviewRemoteReverseZonesModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworkviewRemoteReverseZonesModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworkviewRemoteReverseZones {
	if m == nil {
		return nil
	}
	to := &ipam.NetworkviewRemoteReverseZones{
		Fqdn:                flex.ExpandStringPointer(m.Fqdn),
		ServerAddress:       flex.ExpandIPAddress(m.ServerAddress),
		GssTsigDnsPrincipal: flex.ExpandStringPointer(m.GssTsigDnsPrincipal),
		GssTsigDomain:       flex.ExpandStringPointer(m.GssTsigDomain),
		TsigKey:             flex.ExpandStringPointer(m.TsigKey),
		TsigKeyAlg:          flex.ExpandStringPointer(m.TsigKeyAlg),
		TsigKeyName:         flex.ExpandStringPointer(m.TsigKeyName),
		KeyType:             flex.ExpandStringPointer(m.KeyType),
	}
	return to
}

func FlattenNetworkviewRemoteReverseZones(ctx context.Context, from *ipam.NetworkviewRemoteReverseZones, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworkviewRemoteReverseZonesAttrTypes)
	}
	m := NetworkviewRemoteReverseZonesModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworkviewRemoteReverseZonesAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworkviewRemoteReverseZonesModel) Flatten(ctx context.Context, from *ipam.NetworkviewRemoteReverseZones, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworkviewRemoteReverseZonesModel{}
	}
	m.Fqdn = flex.FlattenStringPointer(from.Fqdn)
	m.ServerAddress = flex.FlattenIPAddress(from.ServerAddress)
	m.GssTsigDnsPrincipal = flex.FlattenStringPointer(from.GssTsigDnsPrincipal)
	m.GssTsigDomain = flex.FlattenStringPointer(from.GssTsigDomain)
	m.TsigKey = flex.FlattenStringPointer(from.TsigKey)
	m.TsigKeyAlg = flex.FlattenStringPointer(from.TsigKeyAlg)
	m.TsigKeyName = flex.FlattenStringPointer(from.TsigKeyName)
	m.KeyType = flex.FlattenStringPointer(from.KeyType)
}

func (m *NetworkviewRemoteReverseZonesModel) PutExpand(to *ipam.NetworkviewRemoteReverseZones) *ipam.NetworkviewRemoteReverseZones {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NetworkviewRemoteReverseZonesResourceSchemaAttributes {
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
