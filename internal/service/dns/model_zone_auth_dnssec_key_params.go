package dns

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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneAuthDnssecKeyParamsModel struct {
	EnableKskAutoRollover         types.Bool   `tfsdk:"enable_ksk_auto_rollover"`
	KskAlgorithm                  types.String `tfsdk:"ksk_algorithm"`
	KskAlgorithms                 types.List   `tfsdk:"ksk_algorithms"`
	KskRollover                   types.Int64  `tfsdk:"ksk_rollover"`
	KskSize                       types.Int64  `tfsdk:"ksk_size"`
	NextSecureType                types.String `tfsdk:"next_secure_type"`
	KskRolloverNotificationConfig types.String `tfsdk:"ksk_rollover_notification_config"`
	KskSnmpNotificationEnabled    types.Bool   `tfsdk:"ksk_snmp_notification_enabled"`
	KskEmailNotificationEnabled   types.Bool   `tfsdk:"ksk_email_notification_enabled"`
	Nsec3SaltMinLength            types.Int64  `tfsdk:"nsec3_salt_min_length"`
	Nsec3SaltMaxLength            types.Int64  `tfsdk:"nsec3_salt_max_length"`
	Nsec3Iterations               types.Int64  `tfsdk:"nsec3_iterations"`
	SignatureExpiration           types.Int64  `tfsdk:"signature_expiration"`
	ZskAlgorithm                  types.String `tfsdk:"zsk_algorithm"`
	ZskAlgorithms                 types.List   `tfsdk:"zsk_algorithms"`
	ZskRollover                   types.Int64  `tfsdk:"zsk_rollover"`
	ZskRolloverMechanism          types.String `tfsdk:"zsk_rollover_mechanism"`
	ZskSize                       types.Int64  `tfsdk:"zsk_size"`
}

var ZoneAuthDnssecKeyParamsAttrTypes = map[string]attr.Type{
	"enable_ksk_auto_rollover":         types.BoolType,
	"ksk_algorithm":                    types.StringType,
	"ksk_algorithms":                   types.ListType{ElemType: types.ObjectType{AttrTypes: ZoneauthdnsseckeyparamsKskAlgorithmsAttrTypes}},
	"ksk_rollover":                     types.Int64Type,
	"ksk_size":                         types.Int64Type,
	"next_secure_type":                 types.StringType,
	"ksk_rollover_notification_config": types.StringType,
	"ksk_snmp_notification_enabled":    types.BoolType,
	"ksk_email_notification_enabled":   types.BoolType,
	"nsec3_salt_min_length":            types.Int64Type,
	"nsec3_salt_max_length":            types.Int64Type,
	"nsec3_iterations":                 types.Int64Type,
	"signature_expiration":             types.Int64Type,
	"zsk_algorithm":                    types.StringType,
	"zsk_algorithms":                   types.ListType{ElemType: types.ObjectType{AttrTypes: ZoneauthdnsseckeyparamsZskAlgorithmsAttrTypes}},
	"zsk_rollover":                     types.Int64Type,
	"zsk_rollover_mechanism":           types.StringType,
	"zsk_size":                         types.Int64Type,
}

var ZoneAuthDnssecKeyParamsResourceSchemaAttributes = map[string]schema.Attribute{
	"enable_ksk_auto_rollover": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If set to True, automatic rollovers for the signing key is enabled.",
	},
	"ksk_algorithm": schema.StringAttribute{
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("10", "13", "14", "5", "7", "8"),
		},
		MarkdownDescription: "Key Signing Key algorithm. Use KSK Algorithms instead.",
	},
	"ksk_algorithms": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZoneauthdnsseckeyparamsKskAlgorithmsResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of Key Signing Key Algorithms.",
	},
	"ksk_rollover": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Key Signing Key rollover interval, in seconds.",
	},
	"ksk_size": schema.Int64Attribute{
		Computed:            true,
		MarkdownDescription: "Key Signing Key size, in bits. Use KSK Algorithms instead.",
	},
	"next_secure_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("NSEC", "NSEC3"),
		},
		MarkdownDescription: "NSEC (next secure) types.",
	},
	"ksk_rollover_notification_config": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ALL", "NONE", "REQUIRE_MANUAL_INTERVENTION", "AUTOMATIC"),
		},
		MarkdownDescription: "This field controls events for which users will be notified.",
	},
	"ksk_snmp_notification_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Enable SNMP notifications for KSK related events.",
	},
	"ksk_email_notification_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Enable email notifications for KSK related events.",
	},
	"nsec3_salt_min_length": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The minimum length for NSEC3 salts.",
	},
	"nsec3_salt_max_length": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The maximum length for NSEC3 salts.",
	},
	"nsec3_iterations": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The number of iterations used for hashing NSEC3.",
	},
	"signature_expiration": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Signature expiration time, in seconds.",
	},
	"zsk_algorithm": schema.StringAttribute{
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("10", "13", "14", "5", "7", "8"),
		},
		MarkdownDescription: "Zone Signing Key algorithm. Use ZSK Algorithms instead.",
	},
	"zsk_algorithms": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZoneauthdnsseckeyparamsZskAlgorithmsResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of Zone Signing Key Algorithms.",
	},
	"zsk_rollover": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Zone Signing Key rollover interval, in seconds.",
	},
	"zsk_rollover_mechanism": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("PRE_PUBLISH", "DOUBLE_SIGN"),
		},
		MarkdownDescription: "Zone Signing Key rollover mechanism.",
	},
	"zsk_size": schema.Int64Attribute{
		Computed:            true,
		MarkdownDescription: "Zone Signing Key size, in bits. Use ZSK Algorithms instead.",
	},
}

func ExpandZoneAuthDnssecKeyParams(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthDnssecKeyParams {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthDnssecKeyParamsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthDnssecKeyParamsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthDnssecKeyParams {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthDnssecKeyParams{
		EnableKskAutoRollover:         flex.ExpandBoolPointer(m.EnableKskAutoRollover),
		KskAlgorithm:                  flex.ExpandStringPointer(m.KskAlgorithm),
		KskAlgorithms:                 flex.ExpandFrameworkListNestedBlock(ctx, m.KskAlgorithms, diags, ExpandZoneauthdnsseckeyparamsKskAlgorithms),
		KskRollover:                   flex.ExpandInt64Pointer(m.KskRollover),
		KskSize:                       flex.ExpandInt64Pointer(m.KskSize),
		NextSecureType:                flex.ExpandStringPointer(m.NextSecureType),
		KskRolloverNotificationConfig: flex.ExpandStringPointer(m.KskRolloverNotificationConfig),
		KskSnmpNotificationEnabled:    flex.ExpandBoolPointer(m.KskSnmpNotificationEnabled),
		KskEmailNotificationEnabled:   flex.ExpandBoolPointer(m.KskEmailNotificationEnabled),
		Nsec3SaltMinLength:            flex.ExpandInt64Pointer(m.Nsec3SaltMinLength),
		Nsec3SaltMaxLength:            flex.ExpandInt64Pointer(m.Nsec3SaltMaxLength),
		Nsec3Iterations:               flex.ExpandInt64Pointer(m.Nsec3Iterations),
		SignatureExpiration:           flex.ExpandInt64Pointer(m.SignatureExpiration),
		ZskAlgorithm:                  flex.ExpandStringPointer(m.ZskAlgorithm),
		ZskAlgorithms:                 flex.ExpandFrameworkListNestedBlock(ctx, m.ZskAlgorithms, diags, ExpandZoneauthdnsseckeyparamsZskAlgorithms),
		ZskRollover:                   flex.ExpandInt64Pointer(m.ZskRollover),
		ZskRolloverMechanism:          flex.ExpandStringPointer(m.ZskRolloverMechanism),
		ZskSize:                       flex.ExpandInt64Pointer(m.ZskSize),
	}
	return to
}

func FlattenZoneAuthDnssecKeyParams(ctx context.Context, from *dns.ZoneAuthDnssecKeyParams, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthDnssecKeyParamsAttrTypes)
	}
	m := ZoneAuthDnssecKeyParamsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthDnssecKeyParamsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthDnssecKeyParamsModel) Flatten(ctx context.Context, from *dns.ZoneAuthDnssecKeyParams, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthDnssecKeyParamsModel{}
	}
	m.EnableKskAutoRollover = types.BoolPointerValue(from.EnableKskAutoRollover)
	m.KskAlgorithm = flex.FlattenStringPointer(from.KskAlgorithm)
	m.KskAlgorithms = flex.FlattenFrameworkListNestedBlock(ctx, from.KskAlgorithms, ZoneauthdnsseckeyparamsKskAlgorithmsAttrTypes, diags, FlattenZoneauthdnsseckeyparamsKskAlgorithms)
	m.KskRollover = flex.FlattenInt64Pointer(from.KskRollover)
	m.KskSize = flex.FlattenInt64Pointer(from.KskSize)
	m.NextSecureType = flex.FlattenStringPointer(from.NextSecureType)
	m.KskRolloverNotificationConfig = flex.FlattenStringPointer(from.KskRolloverNotificationConfig)
	m.KskSnmpNotificationEnabled = types.BoolPointerValue(from.KskSnmpNotificationEnabled)
	m.KskEmailNotificationEnabled = types.BoolPointerValue(from.KskEmailNotificationEnabled)
	m.Nsec3SaltMinLength = flex.FlattenInt64Pointer(from.Nsec3SaltMinLength)
	m.Nsec3SaltMaxLength = flex.FlattenInt64Pointer(from.Nsec3SaltMaxLength)
	m.Nsec3Iterations = flex.FlattenInt64Pointer(from.Nsec3Iterations)
	m.SignatureExpiration = flex.FlattenInt64Pointer(from.SignatureExpiration)
	m.ZskAlgorithm = flex.FlattenStringPointer(from.ZskAlgorithm)
	m.ZskAlgorithms = flex.FlattenFrameworkListNestedBlock(ctx, from.ZskAlgorithms, ZoneauthdnsseckeyparamsZskAlgorithmsAttrTypes, diags, FlattenZoneauthdnsseckeyparamsZskAlgorithms)
	m.ZskRollover = flex.FlattenInt64Pointer(from.ZskRollover)
	m.ZskRolloverMechanism = flex.FlattenStringPointer(from.ZskRolloverMechanism)
	m.ZskSize = flex.FlattenInt64Pointer(from.ZskSize)
}

func (m *ZoneAuthDnssecKeyParamsModel) PutExpand(to *dns.ZoneAuthDnssecKeyParams) *dns.ZoneAuthDnssecKeyParams {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthDnssecKeyParamsResourceSchemaAttributes {
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
