package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneAuthScavengingSettingsModel struct {
	EnableScavenging          types.Bool   `tfsdk:"enable_scavenging"`
	EnableRecurrentScavenging types.Bool   `tfsdk:"enable_recurrent_scavenging"`
	EnableAutoReclamation     types.Bool   `tfsdk:"enable_auto_reclamation"`
	EnableRrLastQueried       types.Bool   `tfsdk:"enable_rr_last_queried"`
	EnableZoneLastQueried     types.Bool   `tfsdk:"enable_zone_last_queried"`
	ReclaimAssociatedRecords  types.Bool   `tfsdk:"reclaim_associated_records"`
	ScavengingSchedule        types.Object `tfsdk:"scavenging_schedule"`
	ExpressionList            types.List   `tfsdk:"expression_list"`
	EaExpressionList          types.List   `tfsdk:"ea_expression_list"`
}

var ZoneAuthScavengingSettingsAttrTypes = map[string]attr.Type{
	"enable_scavenging":           types.BoolType,
	"enable_recurrent_scavenging": types.BoolType,
	"enable_auto_reclamation":     types.BoolType,
	"enable_rr_last_queried":      types.BoolType,
	"enable_zone_last_queried":    types.BoolType,
	"reclaim_associated_records":  types.BoolType,
	"scavenging_schedule":         types.ObjectType{AttrTypes: ZoneauthscavengingsettingsScavengingScheduleAttrTypes},
	"expression_list":             types.ListType{ElemType: types.ObjectType{AttrTypes: ZoneauthscavengingsettingsExpressionListAttrTypes}},
	"ea_expression_list":          types.ListType{ElemType: types.ObjectType{AttrTypes: ZoneauthscavengingsettingsEaExpressionListAttrTypes}},
}

var ZoneAuthScavengingSettingsResourceSchemaAttributes = map[string]schema.Attribute{
	"enable_scavenging": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the resource record scavenging is enabled or not.",
	},
	"enable_recurrent_scavenging": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the recurrent resource record scavenging is enabled or not.",
	},
	"enable_auto_reclamation": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the automatic resource record scavenging is enabled or not.",
	},
	"enable_rr_last_queried": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the resource record last queried monitoring in affected zones is enabled or not.",
	},
	"enable_zone_last_queried": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the last queried monitoring for affected zones is enabled or not.",
	},
	"reclaim_associated_records": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "This flag indicates if the associated resource record scavenging is enabled or not.",
	},
	"scavenging_schedule": schema.SingleNestedAttribute{
		Attributes: ZoneauthscavengingsettingsScavengingScheduleResourceSchemaAttributes,
		Optional:   true,
		Computed:   true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The schedule setting for cloud discovery task.",
	},
	"expression_list": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZoneauthscavengingsettingsExpressionListResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The expression list. The particular record is treated as reclaimable if expression condition evaluates to 'true' for given record if scavenging hasn't been manually disabled on a given resource record.",
	},
	"ea_expression_list": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZoneauthscavengingsettingsEaExpressionListResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The extensible attributes expression list. The particular record is treated as reclaimable if extensible attributes expression condition evaluates to 'true' for given record if scavenging hasn't been manually disabled on a given resource record.",
	},
}

func ExpandZoneAuthScavengingSettings(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthScavengingSettings {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthScavengingSettingsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthScavengingSettingsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthScavengingSettings {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthScavengingSettings{
		EnableScavenging:          flex.ExpandBoolPointer(m.EnableScavenging),
		EnableRecurrentScavenging: flex.ExpandBoolPointer(m.EnableRecurrentScavenging),
		EnableAutoReclamation:     flex.ExpandBoolPointer(m.EnableAutoReclamation),
		EnableRrLastQueried:       flex.ExpandBoolPointer(m.EnableRrLastQueried),
		EnableZoneLastQueried:     flex.ExpandBoolPointer(m.EnableZoneLastQueried),
		ReclaimAssociatedRecords:  flex.ExpandBoolPointer(m.ReclaimAssociatedRecords),
		ScavengingSchedule:        ExpandZoneauthscavengingsettingsScavengingSchedule(ctx, m.ScavengingSchedule, diags),
		ExpressionList:            flex.ExpandFrameworkListNestedBlock(ctx, m.ExpressionList, diags, ExpandZoneauthscavengingsettingsExpressionList),
		EaExpressionList:          flex.ExpandFrameworkListNestedBlock(ctx, m.EaExpressionList, diags, ExpandZoneauthscavengingsettingsEaExpressionList),
	}
	return to
}

func FlattenZoneAuthScavengingSettings(ctx context.Context, from *dns.ZoneAuthScavengingSettings, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthScavengingSettingsAttrTypes)
	}
	m := ZoneAuthScavengingSettingsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthScavengingSettingsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthScavengingSettingsModel) Flatten(ctx context.Context, from *dns.ZoneAuthScavengingSettings, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthScavengingSettingsModel{}
	}
	m.EnableScavenging = types.BoolPointerValue(from.EnableScavenging)
	m.EnableRecurrentScavenging = types.BoolPointerValue(from.EnableRecurrentScavenging)
	m.EnableAutoReclamation = types.BoolPointerValue(from.EnableAutoReclamation)
	m.EnableRrLastQueried = types.BoolPointerValue(from.EnableRrLastQueried)
	m.EnableZoneLastQueried = types.BoolPointerValue(from.EnableZoneLastQueried)
	m.ReclaimAssociatedRecords = types.BoolPointerValue(from.ReclaimAssociatedRecords)
	m.ScavengingSchedule = FlattenZoneauthscavengingsettingsScavengingSchedule(ctx, from.ScavengingSchedule, diags)
	m.ExpressionList = flex.FlattenFrameworkListNestedBlock(ctx, from.ExpressionList, ZoneauthscavengingsettingsExpressionListAttrTypes, diags, FlattenZoneauthscavengingsettingsExpressionList)
	m.EaExpressionList = flex.FlattenFrameworkListNestedBlock(ctx, from.EaExpressionList, ZoneauthscavengingsettingsEaExpressionListAttrTypes, diags, FlattenZoneauthscavengingsettingsEaExpressionList)
}

func (m *ZoneAuthScavengingSettingsModel) PutExpand(to *dns.ZoneAuthScavengingSettings) *dns.ZoneAuthScavengingSettings {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthScavengingSettingsResourceSchemaAttributes {
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
							fmt.Printf("Field: %s, Computed: %v, fieldValue: %v, Value: %s\n", field, boolComp, fieldValue, txtFieldValue)
							if ok {
								if !boolComp {
									continue
								} else if txtFieldValue == "" {
									utils.DeleteBy(to, tField.Name)
								}
							} else if txtFieldValue == "" {
								fmt.Printf("Field: %s is marked as computed but is not a bool. Value: %s\n", field, txtFieldValue)
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
