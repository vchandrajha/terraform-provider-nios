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
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberAutomatedTrafficCaptureSettingModel struct {
	TrafficCaptureEnable    types.Bool   `tfsdk:"traffic_capture_enable"`
	Destination             types.String `tfsdk:"destination"`
	Duration                types.Int64  `tfsdk:"duration"`
	IncludeSupportBundle    types.Bool   `tfsdk:"include_support_bundle"`
	KeepLocalCopy           types.Bool   `tfsdk:"keep_local_copy"`
	DestinationHost         types.String `tfsdk:"destination_host"`
	TrafficCaptureDirectory types.String `tfsdk:"traffic_capture_directory"`
	SupportBundleDirectory  types.String `tfsdk:"support_bundle_directory"`
	Username                types.String `tfsdk:"username"`
	Password                types.String `tfsdk:"password"`
}

var MemberAutomatedTrafficCaptureSettingAttrTypes = map[string]attr.Type{
	"traffic_capture_enable":    types.BoolType,
	"destination":               types.StringType,
	"duration":                  types.Int64Type,
	"include_support_bundle":    types.BoolType,
	"keep_local_copy":           types.BoolType,
	"destination_host":          types.StringType,
	"traffic_capture_directory": types.StringType,
	"support_bundle_directory":  types.StringType,
	"username":                  types.StringType,
	"password":                  types.StringType,
}

var MemberAutomatedTrafficCaptureSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"traffic_capture_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable automated traffic capture based on monitoring thresholds.",
	},
	"destination": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("NONE"),
		Validators: []validator.String{
			stringvalidator.OneOf("FTP", "NONE", "SCP"),
		},
		MarkdownDescription: "Destination of traffic capture files. Save traffic capture locally or upload to remote server using FTP or SCP.",
	},
	"duration": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The time interval on which traffic will be captured(in sec).",
	},
	"include_support_bundle": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Enable automatic download for support bundle.",
	},
	"keep_local_copy": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Save traffic capture files locally.",
	},
	"destination_host": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "IP Address of the destination host.",
	},
	"traffic_capture_directory": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Directory to store the traffic capture files on the remote server.",
	},
	"support_bundle_directory": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Directory to store the support bundle on the remote server.",
	},
	"username": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "User name for accessing the FTP/SCP server.",
	},
	"password": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Password for accessing the FTP/SCP server. This field is not readable.",
	},
}

func ExpandMemberAutomatedTrafficCaptureSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberAutomatedTrafficCaptureSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberAutomatedTrafficCaptureSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberAutomatedTrafficCaptureSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberAutomatedTrafficCaptureSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberAutomatedTrafficCaptureSetting{
		TrafficCaptureEnable:    flex.ExpandBoolPointer(m.TrafficCaptureEnable),
		Destination:             flex.ExpandStringPointer(m.Destination),
		Duration:                flex.ExpandInt64Pointer(m.Duration),
		IncludeSupportBundle:    flex.ExpandBoolPointer(m.IncludeSupportBundle),
		KeepLocalCopy:           flex.ExpandBoolPointer(m.KeepLocalCopy),
		DestinationHost:         flex.ExpandStringPointerEmptyAsNil(m.DestinationHost),
		TrafficCaptureDirectory: flex.ExpandStringPointer(m.TrafficCaptureDirectory),
		SupportBundleDirectory:  flex.ExpandStringPointer(m.SupportBundleDirectory),
		Username:                flex.ExpandStringPointer(m.Username),
		Password:                flex.ExpandStringPointer(m.Password),
	}
	return to
}

func FlattenMemberAutomatedTrafficCaptureSetting(ctx context.Context, from *grid.MemberAutomatedTrafficCaptureSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberAutomatedTrafficCaptureSettingAttrTypes)
	}
	m := MemberAutomatedTrafficCaptureSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberAutomatedTrafficCaptureSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberAutomatedTrafficCaptureSettingModel) Flatten(ctx context.Context, from *grid.MemberAutomatedTrafficCaptureSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberAutomatedTrafficCaptureSettingModel{}
	}
	m.TrafficCaptureEnable = types.BoolPointerValue(from.TrafficCaptureEnable)
	m.Destination = flex.FlattenStringPointer(from.Destination)
	m.Duration = flex.FlattenInt64Pointer(from.Duration)
	m.IncludeSupportBundle = types.BoolPointerValue(from.IncludeSupportBundle)
	m.KeepLocalCopy = types.BoolPointerValue(from.KeepLocalCopy)
	m.DestinationHost = flex.FlattenStringPointer(from.DestinationHost)
	m.TrafficCaptureDirectory = flex.FlattenStringPointer(from.TrafficCaptureDirectory)
	m.SupportBundleDirectory = flex.FlattenStringPointer(from.SupportBundleDirectory)
	m.Username = flex.FlattenStringPointer(from.Username)
}

func (m *MemberAutomatedTrafficCaptureSettingModel) PutExpand(to *grid.MemberAutomatedTrafficCaptureSetting) *grid.MemberAutomatedTrafficCaptureSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberAutomatedTrafficCaptureSettingResourceSchemaAttributes {
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
