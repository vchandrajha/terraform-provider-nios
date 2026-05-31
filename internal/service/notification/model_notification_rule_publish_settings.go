package notification

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/notification"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NotificationRulePublishSettingsModel struct {
	EnabledAttributes types.List `tfsdk:"enabled_attributes"`
}

var NotificationRulePublishSettingsAttrTypes = map[string]attr.Type{
	"enabled_attributes": types.ListType{ElemType: types.StringType},
}

var NotificationRulePublishSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"enabled_attributes": schema.ListAttribute{
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.ValueStringsAre(
				stringvalidator.OneOf(
					"CLIENT_ID", "FINGERPRINT", "HOSTNAME", "INFOBLOX_MEMBER", "IPADDRESS",
					"LEASE_END_TIME", "LEASE_START_TIME", "LEASE_STATE", "MAC_OR_DUID", "NETBIOS_NAME",
				),
			),
		},
		Required:            true,
		MarkdownDescription: "The list of NIOS extensible attributes enalbed for publishsing to Cisco ISE endpoint.",
	},
}

func ExpandNotificationRulePublishSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *notification.NotificationRulePublishSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NotificationRulePublishSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NotificationRulePublishSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *notification.NotificationRulePublishSettings {
	if m == nil {
		return nil
	}
	to := &notification.NotificationRulePublishSettings{
		EnabledAttributes: flex.ExpandFrameworkListString(ctx, m.EnabledAttributes, diags),
	}
	return to
}

func FlattenNotificationRulePublishSettings(ctx context.Context, from *notification.NotificationRulePublishSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NotificationRulePublishSettingsAttrTypes)
	}
	m := NotificationRulePublishSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NotificationRulePublishSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NotificationRulePublishSettingsModel) Flatten(ctx context.Context, from *notification.NotificationRulePublishSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NotificationRulePublishSettingsModel{}
	}
	m.EnabledAttributes = flex.FlattenFrameworkListString(ctx, from.EnabledAttributes, diags)
}

func (m *NotificationRulePublishSettingsModel) PutExpand(to *notification.NotificationRulePublishSettings) *notification.NotificationRulePublishSettings {
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

	for field, attr := range NotificationRulePublishSettingsResourceSchemaAttributes {
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
