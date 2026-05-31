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

	for field, attr := range CertificateAuthserviceRemoteLookupServiceResourceSchemaAttributes {
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
