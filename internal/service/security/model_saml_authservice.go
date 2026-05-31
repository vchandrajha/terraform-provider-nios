package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type SamlAuthserviceModel struct {
	Ref            types.String `tfsdk:"ref"`
	Comment        types.String `tfsdk:"comment"`
	Idp            types.Object `tfsdk:"idp"`
	Name           types.String `tfsdk:"name"`
	SessionTimeout types.Int64  `tfsdk:"session_timeout"`
}

var SamlAuthserviceAttrTypes = map[string]attr.Type{
	"ref":             types.StringType,
	"comment":         types.StringType,
	"idp":             types.ObjectType{AttrTypes: SamlAuthserviceIdpAttrTypes},
	"name":            types.StringType,
	"session_timeout": types.Int64Type,
}

var SamlAuthserviceResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The descriptive comment for the SAML authentication service.",
	},
	"idp": schema.SingleNestedAttribute{
		Attributes:          SamlAuthserviceIdpResourceSchemaAttributes,
		Required:            true,
		MarkdownDescription: "The SAML Identity Provider to use for authentication.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of the SAML authentication service.",
	},
	"session_timeout": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(1800),
		MarkdownDescription: "The session timeout in seconds.",
	},
}

func (m *SamlAuthserviceModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.SamlAuthservice {
	if m == nil {
		return nil
	}
	to := &security.SamlAuthservice{
		Comment:        flex.ExpandStringPointer(m.Comment),
		Idp:            ExpandSamlAuthserviceIdp(ctx, m.Idp, diags),
		Name:           flex.ExpandStringPointer(m.Name),
		SessionTimeout: flex.ExpandInt64Pointer(m.SessionTimeout),
	}
	return to
}

func FlattenSamlAuthservice(ctx context.Context, from *security.SamlAuthservice, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SamlAuthserviceAttrTypes)
	}
	m := SamlAuthserviceModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, SamlAuthserviceAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SamlAuthserviceModel) Flatten(ctx context.Context, from *security.SamlAuthservice, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SamlAuthserviceModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	planIdp := m.Idp
	m.Idp = FlattenSamlAuthserviceIdp(ctx, from.Idp, diags)
	if !planIdp.IsUnknown() {
		idpVal, copyDiags := utils.CopyFieldFromPlanToRespObject(ctx, planIdp, m.Idp, "metadata_file_path")
		diags.Append(copyDiags.Errors()...)
		if !copyDiags.HasError() {
			m.Idp = idpVal.(types.Object)
		}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.SessionTimeout = flex.FlattenInt64Pointer(from.SessionTimeout)
}

func (m *SamlAuthserviceModel) PutExpand(to *security.SamlAuthservice) *security.SamlAuthservice {
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

	for field, attr := range SamlAuthserviceResourceSchemaAttributes {
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
