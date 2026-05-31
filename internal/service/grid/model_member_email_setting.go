package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type MemberEmailSettingModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	FromAddress       types.String `tfsdk:"from_address"`
	Address           types.String `tfsdk:"address"`
	RelayEnabled      types.Bool   `tfsdk:"relay_enabled"`
	Relay             types.String `tfsdk:"relay"`
	Password          types.String `tfsdk:"password"`
	Smtps             types.Bool   `tfsdk:"smtps"`
	PortNumber        types.Int64  `tfsdk:"port_number"`
	UseAuthentication types.Bool   `tfsdk:"use_authentication"`
}

var MemberEmailSettingAttrTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"from_address":       types.StringType,
	"address":            types.StringType,
	"relay_enabled":      types.BoolType,
	"relay":              types.StringType,
	"password":           types.StringType,
	"smtps":              types.BoolType,
	"port_number":        types.Int64Type,
	"use_authentication": types.BoolType,
}

var MemberEmailSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if email notification is enabled or not.",
	},
	"from_address": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The email address of a Grid Member for 'from' field in notification.",
	},
	"address": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The notification email address of a Grid member.",
	},
	"relay_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if email relay is enabled or not.",
	},
	"relay": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The relay name or IP address.",
	},
	"password": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Password to validate from address",
	},
	"smtps": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "SMTP over TLS",
	},
	"port_number": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(25),
		MarkdownDescription: "SMTP port number",
	},
	"use_authentication": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable or disable SMTP auth",
	},
}

func ExpandMemberEmailSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberEmailSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberEmailSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberEmailSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberEmailSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberEmailSetting{
		Enabled:           flex.ExpandBoolPointer(m.Enabled),
		FromAddress:       flex.ExpandStringPointer(m.FromAddress),
		Address:           flex.ExpandStringPointer(m.Address),
		RelayEnabled:      flex.ExpandBoolPointer(m.RelayEnabled),
		Relay:             flex.ExpandStringPointer(m.Relay),
		Password:          flex.ExpandStringPointerEmptyAsNil(m.Password),
		Smtps:             flex.ExpandBoolPointer(m.Smtps),
		PortNumber:        flex.ExpandInt64Pointer(m.PortNumber),
		UseAuthentication: flex.ExpandBoolPointer(m.UseAuthentication),
	}
	return to
}

func FlattenMemberEmailSetting(ctx context.Context, from *grid.MemberEmailSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberEmailSettingAttrTypes)
	}
	m := MemberEmailSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberEmailSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberEmailSettingModel) Flatten(ctx context.Context, from *grid.MemberEmailSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberEmailSettingModel{}
	}
	m.Enabled = types.BoolPointerValue(from.Enabled)
	m.FromAddress = flex.FlattenStringPointer(from.FromAddress)
	m.Address = flex.FlattenStringPointer(from.Address)
	m.RelayEnabled = types.BoolPointerValue(from.RelayEnabled)
	m.Relay = flex.FlattenStringPointer(from.Relay)
	m.Smtps = types.BoolPointerValue(from.Smtps)
	m.PortNumber = flex.FlattenInt64Pointer(from.PortNumber)
	m.UseAuthentication = types.BoolPointerValue(from.UseAuthentication)
}

func (m *MemberEmailSettingModel) PutExpand(to *grid.MemberEmailSetting) *grid.MemberEmailSetting {
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

	for field, attr := range MemberEmailSettingResourceSchemaAttributes {
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
