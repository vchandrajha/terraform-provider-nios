package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/security"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type CertificateAuthserviceOcspRespondersModel struct {
	FqdnOrIp            types.String `tfsdk:"fqdn_or_ip"`
	Port                types.Int64  `tfsdk:"port"`
	Comment             types.String `tfsdk:"comment"`
	Disabled            types.Bool   `tfsdk:"disabled"`
	Certificate         types.String `tfsdk:"certificate"`
	CertificateToken    types.String `tfsdk:"certificate_token"`
	CertificateFilePath types.String `tfsdk:"certificate_file_path"`
}

var CertificateAuthserviceOcspRespondersAttrTypes = map[string]attr.Type{
	"fqdn_or_ip":            types.StringType,
	"port":                  types.Int64Type,
	"comment":               types.StringType,
	"disabled":              types.BoolType,
	"certificate":           types.StringType,
	"certificate_token":     types.StringType,
	"certificate_file_path": types.StringType,
}

var CertificateAuthserviceOcspRespondersResourceSchemaAttributes = map[string]schema.Attribute{
	"fqdn_or_ip": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.Any(
				customvalidator.IsValidFQDN(),
				customvalidator.IsValidIPCIDR(),
			),
		},
		MarkdownDescription: "The FQDN (Fully Qualified Domain Name) or IP address of the server.",
	},
	"port": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(80),
		MarkdownDescription: "The port used for connecting.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The descriptive comment for the OCSP authentication responder.",
	},
	"disabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if this OCSP authentication responder is disabled.",
	},
	"certificate": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The reference to the OCSP responder certificate.",
	},
	"certificate_token": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The token returned by the uploadinit function call in object fileop.",
	},
	"certificate_file_path": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The file path to the certificate.",
	},
}

func ExpandCertificateAuthserviceOcspResponders(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.CertificateAuthserviceOcspResponders {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m CertificateAuthserviceOcspRespondersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *CertificateAuthserviceOcspRespondersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.CertificateAuthserviceOcspResponders {
	if m == nil {
		return nil
	}
	to := &security.CertificateAuthserviceOcspResponders{
		FqdnOrIp:         flex.ExpandStringPointer(m.FqdnOrIp),
		Port:             flex.ExpandInt64Pointer(m.Port),
		Comment:          flex.ExpandStringPointer(m.Comment),
		Disabled:         flex.ExpandBoolPointer(m.Disabled),
		CertificateToken: flex.ExpandStringPointer(m.CertificateToken),
	}
	return to
}

func FlattenCertificateAuthserviceOcspResponders(ctx context.Context, from *security.CertificateAuthserviceOcspResponders, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(CertificateAuthserviceOcspRespondersAttrTypes)
	}
	m := CertificateAuthserviceOcspRespondersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, CertificateAuthserviceOcspRespondersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *CertificateAuthserviceOcspRespondersModel) Flatten(ctx context.Context, from *security.CertificateAuthserviceOcspResponders, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = CertificateAuthserviceOcspRespondersModel{}
	}
	m.FqdnOrIp = flex.FlattenStringPointer(from.FqdnOrIp)
	m.Port = flex.FlattenInt64Pointer(from.Port)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Disabled = types.BoolPointerValue(from.Disabled)
	m.CertificateToken = flex.FlattenStringPointer(from.CertificateToken)
}

func (m *CertificateAuthserviceOcspRespondersModel) PutExpand(to *security.CertificateAuthserviceOcspResponders) *security.CertificateAuthserviceOcspResponders {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range CertificateAuthserviceOcspRespondersResourceSchemaAttributes {
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
