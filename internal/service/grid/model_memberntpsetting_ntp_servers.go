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

type MemberntpsettingNtpServersModel struct {
	Address              types.String `tfsdk:"address"`
	EnableAuthentication types.Bool   `tfsdk:"enable_authentication"`
	NtpKeyNumber         types.Int64  `tfsdk:"ntp_key_number"`
	Preferred            types.Bool   `tfsdk:"preferred"`
	Burst                types.Bool   `tfsdk:"burst"`
	Iburst               types.Bool   `tfsdk:"iburst"`
}

var MemberntpsettingNtpServersAttrTypes = map[string]attr.Type{
	"address":               types.StringType,
	"enable_authentication": types.BoolType,
	"ntp_key_number":        types.Int64Type,
	"preferred":             types.BoolType,
	"burst":                 types.BoolType,
	"iburst":                types.BoolType,
}

var MemberntpsettingNtpServersResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The NTP server IP address or FQDN.",
	},
	"enable_authentication": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether the NTP authentication is enabled.",
	},
	"ntp_key_number": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "The NTP authentication key number.",
	},
	"preferred": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether the NTP server is a preferred one or not.",
	},
	"burst": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether the BURST operation mode is enabled. In BURST operating mode, when the external server is reachable and a valid source of synchronization is available, NTP sends a burst of 8 packets with a 2 second interval between packets.",
	},
	"iburst": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Determines whether the IBURST operation mode is enabled. In IBURST operating mode, when the external server is unreachable, NTP server sends a burst of 8 packets with a 2 second interval between packets.",
	},
}

func ExpandMemberntpsettingNtpServers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberntpsettingNtpServers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberntpsettingNtpServersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberntpsettingNtpServersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberntpsettingNtpServers {
	if m == nil {
		return nil
	}
	to := &grid.MemberntpsettingNtpServers{
		Address:              flex.ExpandStringPointer(m.Address),
		EnableAuthentication: flex.ExpandBoolPointer(m.EnableAuthentication),
		NtpKeyNumber:         flex.ExpandInt64Pointer(m.NtpKeyNumber),
		Preferred:            flex.ExpandBoolPointer(m.Preferred),
		Burst:                flex.ExpandBoolPointer(m.Burst),
		Iburst:               flex.ExpandBoolPointer(m.Iburst),
	}
	return to
}

func FlattenMemberntpsettingNtpServers(ctx context.Context, from *grid.MemberntpsettingNtpServers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberntpsettingNtpServersAttrTypes)
	}
	m := MemberntpsettingNtpServersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberntpsettingNtpServersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberntpsettingNtpServersModel) Flatten(ctx context.Context, from *grid.MemberntpsettingNtpServers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberntpsettingNtpServersModel{}
	}
	m.Address = flex.FlattenStringPointer(from.Address)
	m.EnableAuthentication = types.BoolPointerValue(from.EnableAuthentication)
	m.NtpKeyNumber = flex.FlattenInt64Pointer(from.NtpKeyNumber)
	m.Preferred = types.BoolPointerValue(from.Preferred)
	m.Burst = types.BoolPointerValue(from.Burst)
	m.Iburst = types.BoolPointerValue(from.Iburst)
}

func (m *MemberntpsettingNtpServersModel) PutExpand(to *grid.MemberntpsettingNtpServers) *grid.MemberntpsettingNtpServers {
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

	for field, attr := range MemberntpsettingNtpServersResourceSchemaAttributes {
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
