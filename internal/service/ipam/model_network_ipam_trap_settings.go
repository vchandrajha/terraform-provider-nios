package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NetworkIpamTrapSettingsModel struct {
	EnableEmailWarnings types.Bool `tfsdk:"enable_email_warnings"`
	EnableSnmpWarnings  types.Bool `tfsdk:"enable_snmp_warnings"`
}

var NetworkIpamTrapSettingsAttrTypes = map[string]attr.Type{
	"enable_email_warnings": types.BoolType,
	"enable_snmp_warnings":  types.BoolType,
}

var NetworkIpamTrapSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"enable_email_warnings": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether sending warnings by email is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	},
	"enable_snmp_warnings": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether sending warnings by SNMP is enabled or not.",
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
}

func ExpandNetworkIpamTrapSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworkIpamTrapSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworkIpamTrapSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworkIpamTrapSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworkIpamTrapSettings {
	if m == nil {
		return nil
	}
	to := &ipam.NetworkIpamTrapSettings{
		EnableEmailWarnings: flex.ExpandBoolPointer(m.EnableEmailWarnings),
		EnableSnmpWarnings:  flex.ExpandBoolPointer(m.EnableSnmpWarnings),
	}
	return to
}

func FlattenNetworkIpamTrapSettings(ctx context.Context, from *ipam.NetworkIpamTrapSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworkIpamTrapSettingsAttrTypes)
	}
	m := NetworkIpamTrapSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworkIpamTrapSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworkIpamTrapSettingsModel) Flatten(ctx context.Context, from *ipam.NetworkIpamTrapSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworkIpamTrapSettingsModel{}
	}
	m.EnableEmailWarnings = types.BoolPointerValue(from.EnableEmailWarnings)
	m.EnableSnmpWarnings = types.BoolPointerValue(from.EnableSnmpWarnings)
}

func (m *NetworkIpamTrapSettingsModel) PutExpand(to *ipam.NetworkIpamTrapSettings) *ipam.NetworkIpamTrapSettings {
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

	for field, attr := range NetworkIpamTrapSettingsResourceSchemaAttributes {
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
