package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type CertificateAuthserviceRemoteLookupServiceModel struct {
	Ref types.String `tfsdk:"ref"`
}

var CertificateAuthserviceRemoteLookupServiceAttrTypes = map[string]attr.Type{
	"ref": types.StringType,
}

var CertificateAuthserviceRemoteLookupServiceResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The reference to the LDAP auth service object.",
	},
}

func ExpandCertificateAuthserviceRemoteLookupService(ctx context.Context, s types.String, diags *diag.Diagnostics) *security.CertificateAuthserviceRemoteLookupService {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}

	stringPtr := flex.ExpandStringPointer(s)
	if stringPtr == nil {
		return nil
	}

	return &security.CertificateAuthserviceRemoteLookupService{
		String: stringPtr,
	}
}

func (m *CertificateAuthserviceRemoteLookupServiceModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.CertificateAuthserviceRemoteLookupService {
	if m == nil {
		return nil
	}
	to := &security.CertificateAuthserviceRemoteLookupService{
		String: flex.ExpandStringPointer(m.Ref),
	}
	return to
}

func FlattenCertificateAuthserviceRemoteLookupService(ctx context.Context, from *security.CertificateAuthserviceRemoteLookupService, diags *diag.Diagnostics) types.String {
	if from == nil {
		return types.StringNull()
	}
	if from.CertificateAuthserviceRemoteLookupServiceOneOf == nil || from.CertificateAuthserviceRemoteLookupServiceOneOf.Ref == nil {
		return types.StringNull()
	}
	t := from.CertificateAuthserviceRemoteLookupServiceOneOf.Ref
	return flex.FlattenStringPointer(t)
}

func (m *CertificateAuthserviceRemoteLookupServiceModel) Flatten(ctx context.Context, from *security.CertificateAuthserviceRemoteLookupService, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = CertificateAuthserviceRemoteLookupServiceModel{}
	}
	// Check if OneOf structure exists
	if from.CertificateAuthserviceRemoteLookupServiceOneOf == nil || from.CertificateAuthserviceRemoteLookupServiceOneOf.Ref == nil {
		m.Ref = types.StringNull()
		return
	}

	// Extract the Ref from the OneOf structure
	m.Ref = flex.FlattenStringPointer(from.CertificateAuthserviceRemoteLookupServiceOneOf.Ref)
}

func (m *CertificateAuthserviceRemoteLookupServiceModel) PutExpand(to *security.CertificateAuthserviceRemoteLookupService) *security.CertificateAuthserviceRemoteLookupService {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range CertificateAuthserviceRemoteLookupServiceResourceSchemaAttributes {
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
