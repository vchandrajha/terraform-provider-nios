package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

type MemberLomUsersModel struct {
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	Role     types.String `tfsdk:"role"`
	Disable  types.Bool   `tfsdk:"disable"`
	Comment  types.String `tfsdk:"comment"`
}

var MemberLomUsersAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"password": types.StringType,
	"role":     types.StringType,
	"disable":  types.BoolType,
	"comment":  types.StringType,
}

var MemberLomUsersResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The LOM user name.",
	},
	"password": schema.StringAttribute{
		Required:  true,
		Sensitive: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The LOM user password.",
	},
	"role": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("USER"),
		Validators: []validator.String{
			stringvalidator.OneOf("OPERATOR", "USER"),
		},
		MarkdownDescription: "The LOM user role which specifies the list of actions that are allowed for the user.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines whether the LOM user is disabled.",
	},
	"comment": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The descriptive comment for the LOM user.",
	},
}

func ExpandMemberLomUsers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberLomUsers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberLomUsersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags, true)
}

func (m *MemberLomUsersModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *grid.MemberLomUsers {
	if m == nil {
		return nil
	}
	to := &grid.MemberLomUsers{
		Name:     flex.ExpandStringPointer(m.Name),
		Role:     flex.ExpandStringPointer(m.Role),
		Disable:  flex.ExpandBoolPointer(m.Disable),
		Comment:  flex.ExpandStringPointer(m.Comment),
		Password: flex.ExpandStringPointer(m.Password),
	}
	return to
}

func FlattenMemberLomUsers(ctx context.Context, from *grid.MemberLomUsers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberLomUsersAttrTypes)
	}
	m := MemberLomUsersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberLomUsersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberLomUsersModel) Flatten(ctx context.Context, from *grid.MemberLomUsers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberLomUsersModel{}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Password = flex.FlattenStringPointer(from.Password)
	m.Role = flex.FlattenStringPointer(from.Role)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.Comment = flex.FlattenStringPointer(from.Comment)
}

func (m *MemberLomUsersModel) PutExpand(to *grid.MemberLomUsers) *grid.MemberLomUsers {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberLomUsersResourceSchemaAttributes {
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
