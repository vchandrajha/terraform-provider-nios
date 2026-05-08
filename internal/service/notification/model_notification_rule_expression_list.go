package notification

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/notification"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type NotificationRuleExpressionListModel struct {
	Op      types.String `tfsdk:"op"`
	Op1     types.String `tfsdk:"op1"`
	Op1Type types.String `tfsdk:"op1_type"`
	Op2     types.String `tfsdk:"op2"`
	Op2Type types.String `tfsdk:"op2_type"`
}

var NotificationRuleExpressionListAttrTypes = map[string]attr.Type{
	"op":       types.StringType,
	"op1":      types.StringType,
	"op1_type": types.StringType,
	"op2":      types.StringType,
	"op2_type": types.StringType,
}

var NotificationRuleExpressionListResourceSchemaAttributes = map[string]schema.Attribute{
	"op": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("AND", "CONTAINED_IN", "ENDLIST", "EQ", "EXISTS",
				"GE", "GT", "LE", "LT", "MATCH_CIDR", "MATCH_RANGE", "NOT_EQ",
				"NOT_EXISTS", "NREGEX", "OR", "REGEX"),
		},
		MarkdownDescription: "Rule expression type.",
	},
	"op1": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf(
				"ADDRESS_TYPE", "ATC_HIT_CLASS", "ATC_HIT_PROPERTY", "ATC_HIT_TYPE",
				"AUTO_CREATED_RECORDS", "DB_CHANGE_GROUP_NAME", "DB_CHANGE_USER_NAME",
				"DHCP_FINGERPRINT", "DHCP_IP_ADDRESS", "DHCP_LEASE_STATE", "DISABLE",
				"DISCOVERER", "DNS_RPZ_ACTION_POLICY", "DNS_RPZ_NAME", "DNS_RPZ_RULE_NAME",
				"DNS_RPZ_TYPE", "DNS_VIEW", "DUID", "DXL_TOPIC", "HOST", "IPV4_ADDRESS",
				"IPV6_ADDRESS", "IPV6_PREFIX", "IPV6_PREFIX_BITS", "IP_ADDRESS", "IS_IPV4",
				"MAC", "MEMBER_IP", "MEMBER_NAME", "NAME", "NETWORK", "NETWORK_VIEW",
				"OPERATION_TYPE", "QUERY_FQDN", "RECORD_NAME", "RECORD_TYPE",
				"SECURITY_ADP_HITS_COUNT", "SECURITY_ADP_RULE_ACTION",
				"SECURITY_ADP_RULE_CATEGORY", "SECURITY_ADP_RULE_MESSAGE",
				"SECURITY_ADP_RULE_SEVERITY", "SECURITY_ADP_SID", "SERVER_ASSOC_TYPE_IPV4",
				"SERVER_ASSOC_TYPE_IPV6", "SOURCE_IP", "THREAT_ORIGIN", "UNMANAGED", "ZONE_NAME",
				"ZONE_TYPE",
			),
		},
		MarkdownDescription: "Rule expression first operand value.",
	},
	"op1_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("FIELD", "LIST", "STRING"),
		},
		MarkdownDescription: "Rule expression first operand type.",
	},
	"op2": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Rule expression second operand.",
	},
	"op2_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("FIELD", "LIST", "STRING"),
		},
		MarkdownDescription: "Rule expression second operand type.",
	},
}

func ExpandNotificationRuleExpressionList(ctx context.Context, o types.Object, diags *diag.Diagnostics) *notification.NotificationRuleExpressionList {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NotificationRuleExpressionListModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NotificationRuleExpressionListModel) Expand(ctx context.Context, diags *diag.Diagnostics) *notification.NotificationRuleExpressionList {
	if m == nil {
		return nil
	}
	to := &notification.NotificationRuleExpressionList{
		Op:      flex.ExpandStringPointer(m.Op),
		Op1:     flex.ExpandStringPointer(m.Op1),
		Op1Type: flex.ExpandStringPointer(m.Op1Type),
		Op2:     flex.ExpandStringPointer(m.Op2),
		Op2Type: flex.ExpandStringPointer(m.Op2Type),
	}
	return to
}

func FlattenNotificationRuleExpressionList(ctx context.Context, from *notification.NotificationRuleExpressionList, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NotificationRuleExpressionListAttrTypes)
	}
	m := NotificationRuleExpressionListModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NotificationRuleExpressionListAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NotificationRuleExpressionListModel) Flatten(ctx context.Context, from *notification.NotificationRuleExpressionList, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NotificationRuleExpressionListModel{}
	}
	m.Op = flex.FlattenStringPointer(from.Op)
	m.Op1 = flex.FlattenStringPointer(from.Op1)
	m.Op1Type = flex.FlattenStringPointer(from.Op1Type)
	m.Op2 = flex.FlattenStringPointer(from.Op2)
	m.Op2Type = flex.FlattenStringPointer(from.Op2Type)
}

func (m *NotificationRuleExpressionListModel) PutExpand(to *notification.NotificationRuleExpressionList) *notification.NotificationRuleExpressionList {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NotificationRuleExpressionListResourceSchemaAttributes {
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
							fmt.Printf("Field: %s, ok: %v, Computed: %v, fieldValue: %v, Value: %s\n", field, ok, boolComp, fieldValue, txtFieldValue)
							if ok {
								if boolComp && txtFieldValue == "" {
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
