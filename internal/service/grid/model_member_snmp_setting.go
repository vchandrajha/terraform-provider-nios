package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type MemberSnmpSettingModel struct {
	EngineId               types.List   `tfsdk:"engine_id"`
	QueriesCommunityString types.String `tfsdk:"queries_community_string"`
	QueriesEnable          types.Bool   `tfsdk:"queries_enable"`
	Snmpv3QueriesEnable    types.Bool   `tfsdk:"snmpv3_queries_enable"`
	Snmpv3QueriesUsers     types.List   `tfsdk:"snmpv3_queries_users"`
	Snmpv3TrapsEnable      types.Bool   `tfsdk:"snmpv3_traps_enable"`
	Syscontact             types.List   `tfsdk:"syscontact"`
	Sysdescr               types.List   `tfsdk:"sysdescr"`
	Syslocation            types.List   `tfsdk:"syslocation"`
	Sysname                types.List   `tfsdk:"sysname"`
	TrapReceivers          types.List   `tfsdk:"trap_receivers"`
	TrapsCommunityString   types.String `tfsdk:"traps_community_string"`
	TrapsEnable            types.Bool   `tfsdk:"traps_enable"`
}

var MemberSnmpSettingAttrTypes = map[string]attr.Type{
	"engine_id":                types.ListType{ElemType: types.StringType},
	"queries_community_string": types.StringType,
	"queries_enable":           types.BoolType,
	"snmpv3_queries_enable":    types.BoolType,
	"snmpv3_queries_users":     types.ListType{ElemType: types.ObjectType{AttrTypes: MembersnmpsettingSnmpv3QueriesUsersAttrTypes}},
	"snmpv3_traps_enable":      types.BoolType,
	"syscontact":               types.ListType{ElemType: types.StringType},
	"sysdescr":                 types.ListType{ElemType: types.StringType},
	"syslocation":              types.ListType{ElemType: types.StringType},
	"sysname":                  types.ListType{ElemType: types.StringType},
	"trap_receivers":           types.ListType{ElemType: types.ObjectType{AttrTypes: MembersnmpsettingTrapReceiversAttrTypes}},
	"traps_community_string":   types.StringType,
	"traps_enable":             types.BoolType,
}

var MemberSnmpSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"engine_id": schema.ListAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The engine ID of the appliance that manages the SNMP agent.",
	},
	"queries_community_string": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The community string for SNMP queries.",
	},
	"queries_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If set to True, SNMP queries are enabled.",
	},
	"snmpv3_queries_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If set to True, SNMPv3 queries are enabled.",
	},
	"snmpv3_queries_users": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MembersnmpsettingSnmpv3QueriesUsersResourceSchemaAttributes,
		},
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of SNMPv3 queries users.",
	},
	"snmpv3_traps_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If set to True, SNMPv3 traps are enabled.",
	},
	"syscontact": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The name of the contact person for the appliance. Second value is applicable only for HA pair. Otherwise second value is ignored.",
	},
	"sysdescr": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Useful information about the appliance. Second value is applicable only for HA pair. Otherwise second value is ignored.",
	},
	"syslocation": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The physical location of the appliance. Second value is applicable only for HA pair. Otherwise second value is ignored.",
	},
	"sysname": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The FQDN (Fully Qualified Domain Name) of the appliance. Second value is applicable only for HA pair. Otherwise second value is ignored.",
	},
	"trap_receivers": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MembersnmpsettingTrapReceiversResourceSchemaAttributes,
		},
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of trap receivers.",
	},
	"traps_community_string": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "A string the NIOS appliance sends to the management system together with its traps. Note that this community string must match exactly what you enter in the management system.",
	},
	"traps_enable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If set to True, SNMP traps are enabled.",
	},
}

func ExpandMemberSnmpSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberSnmpSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberSnmpSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberSnmpSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberSnmpSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberSnmpSetting{
		QueriesCommunityString: flex.ExpandStringPointer(m.QueriesCommunityString),
		QueriesEnable:          flex.ExpandBoolPointer(m.QueriesEnable),
		Snmpv3QueriesEnable:    flex.ExpandBoolPointer(m.Snmpv3QueriesEnable),
		Snmpv3QueriesUsers:     flex.ExpandFrameworkListNestedBlock(ctx, m.Snmpv3QueriesUsers, diags, ExpandMembersnmpsettingSnmpv3QueriesUsers),
		Snmpv3TrapsEnable:      flex.ExpandBoolPointer(m.Snmpv3TrapsEnable),
		Syscontact:             flex.ExpandFrameworkListString(ctx, m.Syscontact, diags),
		Sysdescr:               flex.ExpandFrameworkListString(ctx, m.Sysdescr, diags),
		Syslocation:            flex.ExpandFrameworkListString(ctx, m.Syslocation, diags),
		Sysname:                flex.ExpandFrameworkListString(ctx, m.Sysname, diags),
		TrapReceivers:          flex.ExpandFrameworkListNestedBlock(ctx, m.TrapReceivers, diags, ExpandMembersnmpsettingTrapReceivers),
		TrapsCommunityString:   flex.ExpandStringPointer(m.TrapsCommunityString),
		TrapsEnable:            flex.ExpandBoolPointer(m.TrapsEnable),
	}
	return to
}

func FlattenMemberSnmpSetting(ctx context.Context, from *grid.MemberSnmpSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberSnmpSettingAttrTypes)
	}
	m := MemberSnmpSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberSnmpSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberSnmpSettingModel) Flatten(ctx context.Context, from *grid.MemberSnmpSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberSnmpSettingModel{}
	}
	m.EngineId = flex.FlattenFrameworkListString(ctx, from.EngineId, diags)
	m.QueriesCommunityString = flex.FlattenStringPointer(from.QueriesCommunityString)
	m.QueriesEnable = types.BoolPointerValue(from.QueriesEnable)
	m.Snmpv3QueriesEnable = types.BoolPointerValue(from.Snmpv3QueriesEnable)
	m.Snmpv3QueriesUsers = flex.FlattenFrameworkListNestedBlock(ctx, from.Snmpv3QueriesUsers, MembersnmpsettingSnmpv3QueriesUsersAttrTypes, diags, FlattenMembersnmpsettingSnmpv3QueriesUsers)
	m.Snmpv3TrapsEnable = types.BoolPointerValue(from.Snmpv3TrapsEnable)
	m.Syscontact = flex.FlattenFrameworkListString(ctx, from.Syscontact, diags)
	m.Sysdescr = flex.FlattenFrameworkListString(ctx, from.Sysdescr, diags)
	m.Syslocation = flex.FlattenFrameworkListString(ctx, from.Syslocation, diags)
	m.Sysname = flex.FlattenFrameworkListString(ctx, from.Sysname, diags)
	m.TrapReceivers = flex.FlattenFrameworkListNestedBlock(ctx, from.TrapReceivers, MembersnmpsettingTrapReceiversAttrTypes, diags, FlattenMembersnmpsettingTrapReceivers)
	m.TrapsCommunityString = flex.FlattenStringPointer(from.TrapsCommunityString)
	m.TrapsEnable = types.BoolPointerValue(from.TrapsEnable)
}

func (m *MemberSnmpSettingModel) PutExpand(to *grid.MemberSnmpSetting) *grid.MemberSnmpSetting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberSnmpSettingResourceSchemaAttributes {
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
