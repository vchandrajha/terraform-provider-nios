package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

type LdapAuthServiceModel struct {
	Ref                         types.String `tfsdk:"ref"`
	Comment                     types.String `tfsdk:"comment"`
	Disable                     types.Bool   `tfsdk:"disable"`
	EaMapping                   types.List   `tfsdk:"ea_mapping"`
	LdapGroupAttribute          types.String `tfsdk:"ldap_group_attribute"`
	LdapGroupAuthenticationType types.String `tfsdk:"ldap_group_authentication_type"`
	LdapUserAttribute           types.String `tfsdk:"ldap_user_attribute"`
	Mode                        types.String `tfsdk:"mode"`
	Name                        types.String `tfsdk:"name"`
	RecoveryInterval            types.Int64  `tfsdk:"recovery_interval"`
	Retries                     types.Int64  `tfsdk:"retries"`
	SearchScope                 types.String `tfsdk:"search_scope"`
	Servers                     types.List   `tfsdk:"servers"`
	Timeout                     types.Int64  `tfsdk:"timeout"`
}

var LdapAuthServiceAttrTypes = map[string]attr.Type{
	"ref":                            types.StringType,
	"comment":                        types.StringType,
	"disable":                        types.BoolType,
	"ea_mapping":                     types.ListType{ElemType: types.ObjectType{AttrTypes: LdapAuthServiceEaMappingAttrTypes}},
	"ldap_group_attribute":           types.StringType,
	"ldap_group_authentication_type": types.StringType,
	"ldap_user_attribute":            types.StringType,
	"mode":                           types.StringType,
	"name":                           types.StringType,
	"recovery_interval":              types.Int64Type,
	"retries":                        types.Int64Type,
	"search_scope":                   types.StringType,
	"servers":                        types.ListType{ElemType: types.ObjectType{AttrTypes: LdapAuthServiceServersAttrTypes}},
	"timeout":                        types.Int64Type,
}

var LdapAuthServiceResourceSchemaAttributes = map[string]schema.Attribute{
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
			stringvalidator.LengthBetween(0, 256),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The LDAP descriptive comment.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the LDAP authentication service is disabled.",
	},
	"ea_mapping": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: LdapAuthServiceEaMappingResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Optional:            true,
		MarkdownDescription: "The mapping LDAP fields to extensible attributes.",
	},
	"ldap_group_attribute": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("memberOf"),
		MarkdownDescription: "The name of the LDAP attribute that defines group membership.",
	},
	"ldap_group_authentication_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("GROUP_ATTRIBUTE", "POSIX_GROUP"),
		},
		Default:             stringdefault.StaticString("GROUP_ATTRIBUTE"),
		MarkdownDescription: "The LDAP group authentication type.",
	},
	"ldap_user_attribute": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The LDAP userid attribute that is used for search.",
	},
	"mode": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ORDERED_LIST", "ROUND_ROBIN"),
		},
		Default:             stringdefault.StaticString("ORDERED_LIST"),
		MarkdownDescription: "The LDAP authentication mode.",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The LDAP authentication service name.",
	},
	"recovery_interval": schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 600),
		},
		MarkdownDescription: "The period of time in seconds to wait before trying to contact a LDAP server that has been marked as 'DOWN'.",
	},
	"retries": schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 5),
		},
		MarkdownDescription: "The maximum number of LDAP authentication attempts.",
	},
	"search_scope": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("BASE", "ONELEVEL", "SUBTREE"),
		},
		Default:             stringdefault.StaticString("ONELEVEL"),
		MarkdownDescription: "The starting point of the LDAP search.",
	},
	"servers": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: LdapAuthServiceServersResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Required:            true,
		MarkdownDescription: "The list of LDAP servers used for authentication.",
	},
	"timeout": schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 600),
		},
		MarkdownDescription: "The LDAP authentication timeout in seconds.",
	},
}

func (m *LdapAuthServiceModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.LdapAuthService {
	if m == nil {
		return nil
	}
	to := &security.LdapAuthService{
		Comment:                     flex.ExpandStringPointer(m.Comment),
		Disable:                     flex.ExpandBoolPointer(m.Disable),
		EaMapping:                   flex.ExpandFrameworkListNestedBlock(ctx, m.EaMapping, diags, ExpandLdapAuthServiceEaMapping),
		LdapGroupAttribute:          flex.ExpandStringPointer(m.LdapGroupAttribute),
		LdapGroupAuthenticationType: flex.ExpandStringPointer(m.LdapGroupAuthenticationType),
		LdapUserAttribute:           flex.ExpandStringPointer(m.LdapUserAttribute),
		Mode:                        flex.ExpandStringPointer(m.Mode),
		Name:                        flex.ExpandStringPointer(m.Name),
		RecoveryInterval:            flex.ExpandInt64Pointer(m.RecoveryInterval),
		Retries:                     flex.ExpandInt64Pointer(m.Retries),
		SearchScope:                 flex.ExpandStringPointer(m.SearchScope),
		Servers:                     flex.ExpandFrameworkListNestedBlock(ctx, m.Servers, diags, ExpandLdapAuthServiceServers),
		Timeout:                     flex.ExpandInt64Pointer(m.Timeout),
	}
	return to
}

func FlattenLdapAuthService(ctx context.Context, from *security.LdapAuthService, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(LdapAuthServiceAttrTypes)
	}
	m := LdapAuthServiceModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, LdapAuthServiceAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *LdapAuthServiceModel) Flatten(ctx context.Context, from *security.LdapAuthService, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = LdapAuthServiceModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.EaMapping = flex.FlattenFrameworkListNestedBlock(ctx, from.EaMapping, LdapAuthServiceEaMappingAttrTypes, diags, FlattenLdapAuthServiceEaMapping)
	m.LdapGroupAttribute = flex.FlattenStringPointer(from.LdapGroupAttribute)
	m.LdapGroupAuthenticationType = flex.FlattenStringPointer(from.LdapGroupAuthenticationType)
	m.LdapUserAttribute = flex.FlattenStringPointer(from.LdapUserAttribute)
	m.Mode = flex.FlattenStringPointer(from.Mode)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.RecoveryInterval = flex.FlattenInt64Pointer(from.RecoveryInterval)
	m.Retries = flex.FlattenInt64Pointer(from.Retries)
	m.SearchScope = flex.FlattenStringPointer(from.SearchScope)
	m.Servers = flex.FlattenFrameworkListNestedBlock(ctx, from.Servers, LdapAuthServiceServersAttrTypes, diags, FlattenLdapAuthServiceServers)
	m.Timeout = flex.FlattenInt64Pointer(from.Timeout)
}

func (m *LdapAuthServiceModel) PutExpand(to *security.LdapAuthService) *security.LdapAuthService {
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

	for field, attr := range LdapAuthServiceResourceSchemaAttributes {
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
