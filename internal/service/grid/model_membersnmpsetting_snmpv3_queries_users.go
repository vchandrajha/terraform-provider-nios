package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type MembersnmpsettingSnmpv3QueriesUsersModel struct {
	Ref                    types.String `tfsdk:"ref"`
	User                   types.String `tfsdk:"user"`
	AuthenticationProtocol types.String `tfsdk:"authentication_protocol"`
	Comment                types.String `tfsdk:"comment"`
	Disable                types.Bool   `tfsdk:"disable"`
	ExtAttrs               types.Map    `tfsdk:"extattrs"`
	Name                   types.String `tfsdk:"name"`
	PrivacyProtocol        types.String `tfsdk:"privacy_protocol"`
}

var MembersnmpsettingSnmpv3QueriesUsersAttrTypes = map[string]attr.Type{
	"ref":                     types.StringType,
	"user":                    types.StringType,
	"authentication_protocol": types.StringType,
	"comment":                 types.StringType,
	"disable":                 types.BoolType,
	"extattrs":                types.MapType{ElemType: types.StringType},
	"name":                    types.StringType,
	"privacy_protocol":        types.StringType,
}

var MembersnmpsettingSnmpv3QueriesUsersResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the SNMPv3 user object",
	},
	"user": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The SNMPv3 user.",
	},
	"authentication_protocol": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The authentication protocol to be used for this user.",
	},
	"comment": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "A descriptive comment for the SNMPv3 User.",
	},
	"disable": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if SNMPv3 user is disabled or not.",
	},
	"extattrs": schema.MapAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers: []planmodifier.Map{
			mapplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Extensible attributes associated with the object. For valid values for extensible attributes, see {extattrs:values}.",
	},
	"name": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The name of the user.",
	},
	"privacy_protocol": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The privacy protocol to be used for this user.",
	},
}

func ExpandMembersnmpsettingSnmpv3QueriesUsers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MembersnmpsettingSnmpv3QueriesUsers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MembersnmpsettingSnmpv3QueriesUsersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MembersnmpsettingSnmpv3QueriesUsersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MembersnmpsettingSnmpv3QueriesUsers {
	if m == nil {
		return nil
	}
	var oneOf *grid.MembersnmpsettingSnmpv3QueriesUsersOneOf
	if !m.User.IsNull() && !m.User.IsUnknown() {
		oneOf = &grid.MembersnmpsettingSnmpv3QueriesUsersOneOf{
			User: flex.ExpandStringPointer(m.User),
		}
	}

	to := &grid.MembersnmpsettingSnmpv3QueriesUsers{
		MembersnmpsettingSnmpv3QueriesUsersOneOf: oneOf,
	}

	return to
}

func FlattenMembersnmpsettingSnmpv3QueriesUsers(ctx context.Context, from *grid.MembersnmpsettingSnmpv3QueriesUsers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MembersnmpsettingSnmpv3QueriesUsersAttrTypes)
	}
	m := MembersnmpsettingSnmpv3QueriesUsersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MembersnmpsettingSnmpv3QueriesUsersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MembersnmpsettingSnmpv3QueriesUsersModel) Flatten(ctx context.Context, from *grid.MembersnmpsettingSnmpv3QueriesUsers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MembersnmpsettingSnmpv3QueriesUsersModel{}
	}
	m.User = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.Ref)
	m.Ref = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.Ref)
	m.AuthenticationProtocol = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.AuthenticationProtocol)
	m.Comment = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.Comment)
	m.Disable = types.BoolPointerValue(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.Disable)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.ExtAttrs, diags)
	m.Name = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.Name)
	m.PrivacyProtocol = flex.FlattenStringPointer(from.MembersnmpsettingSnmpv3QueriesUsersOneOf1.User.PrivacyProtocol)
}

func (m *MembersnmpsettingSnmpv3QueriesUsersModel) PutExpand(to *grid.MembersnmpsettingSnmpv3QueriesUsers) *grid.MembersnmpsettingSnmpv3QueriesUsers {
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

	for field, attr := range MembersnmpsettingSnmpv3QueriesUsersResourceSchemaAttributes {
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
