package dhcp

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type FixedaddressCliCredentialsModel struct {
	User            types.String `tfsdk:"user"`
	Password        types.String `tfsdk:"password"`
	CredentialType  types.String `tfsdk:"credential_type"`
	Comment         types.String `tfsdk:"comment"`
	Id              types.Int64  `tfsdk:"id"`
	CredentialGroup types.String `tfsdk:"credential_group"`
}

var FixedaddressCliCredentialsAttrTypes = map[string]attr.Type{
	"user":             types.StringType,
	"password":         types.StringType,
	"credential_type":  types.StringType,
	"comment":          types.StringType,
	"id":               types.Int64Type,
	"credential_group": types.StringType,
}

var FixedaddressCliCredentialsResourceSchemaAttributes = map[string]schema.Attribute{
	"user": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The CLI user name.",
	},
	"password": schema.StringAttribute{
		Optional:  true,
		Computed:  true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Sensitive: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The CLI password.",
	},
	"credential_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ENABLE_SSH", "ENABLE_TELNET", "SSH", "TELNET"),
		},
		MarkdownDescription: "The type of the credential.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The comment for the credential.",
	},
	"id": schema.Int64Attribute{
		Computed:            true,
		MarkdownDescription: "The Credentials ID.",
	},
	"credential_group": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Group for the CLI credential.",
	},
}

func ExpandFixedaddressCliCredentials(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dhcp.FixedaddressCliCredentials {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m FixedaddressCliCredentialsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *FixedaddressCliCredentialsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.FixedaddressCliCredentials {
	if m == nil {
		return nil
	}
	to := &dhcp.FixedaddressCliCredentials{
		User:            flex.ExpandStringPointer(m.User),
		Password:        flex.ExpandStringPointer(m.Password),
		CredentialType:  flex.ExpandStringPointer(m.CredentialType),
		Comment:         flex.ExpandStringPointer(m.Comment),
		CredentialGroup: flex.ExpandStringPointer(m.CredentialGroup),
	}
	return to
}

func FlattenFixedaddressCliCredentials(ctx context.Context, from *dhcp.FixedaddressCliCredentials, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(FixedaddressCliCredentialsAttrTypes)
	}
	m := FixedaddressCliCredentialsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, FixedaddressCliCredentialsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *FixedaddressCliCredentialsModel) Flatten(ctx context.Context, from *dhcp.FixedaddressCliCredentials, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = FixedaddressCliCredentialsModel{}
	}
	m.User = flex.FlattenStringPointer(from.User)
	m.Password = flex.FlattenStringPointer(from.Password)
	m.CredentialType = flex.FlattenStringPointer(from.CredentialType)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Id = flex.FlattenInt64Pointer(from.Id)
	m.CredentialGroup = flex.FlattenStringPointer(from.CredentialGroup)
}

func (m *FixedaddressCliCredentialsModel) PutExpand(to *dhcp.FixedaddressCliCredentials) *dhcp.FixedaddressCliCredentials {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range FixedaddressCliCredentialsResourceSchemaAttributes {
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
