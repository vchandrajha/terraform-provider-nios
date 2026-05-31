package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupLockoutSettingModel struct {
	EnableSequentialFailedLoginAttemptsLockout types.Bool  `tfsdk:"enable_sequential_failed_login_attempts_lockout"`
	SequentialAttempts                         types.Int64 `tfsdk:"sequential_attempts"`
	FailedLockoutDuration                      types.Int64 `tfsdk:"failed_lockout_duration"`
	NeverUnlockUser                            types.Bool  `tfsdk:"never_unlock_user"`
}

var AdmingroupLockoutSettingAttrTypes = map[string]attr.Type{
	"enable_sequential_failed_login_attempts_lockout": types.BoolType,
	"sequential_attempts":                             types.Int64Type,
	"failed_lockout_duration":                         types.Int64Type,
	"never_unlock_user":                               types.BoolType,
}

var AdmingroupLockoutSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"enable_sequential_failed_login_attempts_lockout": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable/disable sequential failed login attempts lockout for local users",
	},
	"sequential_attempts": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(5),
		Validators: []validator.Int64{
			int64validator.Between(1, 99),
		},
		MarkdownDescription: "The number of failed login attempts",
	},
	"failed_lockout_duration": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		Default:  int64default.StaticInt64(5),
		Validators: []validator.Int64{
			int64validator.Between(1, 1440),
		},
		MarkdownDescription: "Time period the account remains locked after sequential failed login attempt lockout.",
	},
	"never_unlock_user": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Never unlock option is also provided and if set then user account is locked forever and only super user can unlock this account",
	},
}

func ExpandAdmingroupLockoutSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupLockoutSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupLockoutSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupLockoutSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupLockoutSetting {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupLockoutSetting{
		EnableSequentialFailedLoginAttemptsLockout: flex.ExpandBoolPointer(m.EnableSequentialFailedLoginAttemptsLockout),
		SequentialAttempts:                         flex.ExpandInt64Pointer(m.SequentialAttempts),
		FailedLockoutDuration:                      flex.ExpandInt64Pointer(m.FailedLockoutDuration),
		NeverUnlockUser:                            flex.ExpandBoolPointer(m.NeverUnlockUser),
	}
	return to
}

func FlattenAdmingroupLockoutSetting(ctx context.Context, from *security.AdmingroupLockoutSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupLockoutSettingAttrTypes)
	}
	m := AdmingroupLockoutSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupLockoutSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupLockoutSettingModel) Flatten(ctx context.Context, from *security.AdmingroupLockoutSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupLockoutSettingModel{}
	}
	m.EnableSequentialFailedLoginAttemptsLockout = types.BoolPointerValue(from.EnableSequentialFailedLoginAttemptsLockout)
	m.SequentialAttempts = flex.FlattenInt64Pointer(from.SequentialAttempts)
	m.FailedLockoutDuration = flex.FlattenInt64Pointer(from.FailedLockoutDuration)
	m.NeverUnlockUser = types.BoolPointerValue(from.NeverUnlockUser)
}

func (m *AdmingroupLockoutSettingModel) PutExpand(to *security.AdmingroupLockoutSetting) *security.AdmingroupLockoutSetting {
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

	for field, attr := range AdmingroupLockoutSettingResourceSchemaAttributes {
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
