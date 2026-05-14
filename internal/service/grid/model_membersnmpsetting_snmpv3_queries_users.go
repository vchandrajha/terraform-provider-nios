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
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
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
		MarkdownDescription: "The reference to the SNMPv3 user object",
	},
	"user": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The SNMPv3 user.",
	},
	"authentication_protocol": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The authentication protocol to be used for this user.",
	},
	"comment": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "A descriptive comment for the SNMPv3 User.",
	},
	"disable": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Determines if SNMPv3 user is disabled or not.",
	},
	"extattrs": schema.MapAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		MarkdownDescription: "Extensible attributes associated with the object. For valid values for extensible attributes, see {extattrs:values}.",
	},
	"name": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The name of the user.",
	},
	"privacy_protocol": schema.StringAttribute{
		Computed:            true,
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
	for field, attr := range MembersnmpsettingSnmpv3QueriesUsersResourceSchemaAttributes {
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
