package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type RecordAAwsRte53RecordInfoModel struct {
	AliasTargetDnsName              types.String `tfsdk:"alias_target_dns_name"`
	AliasTargetHostedZoneId         types.String `tfsdk:"alias_target_hosted_zone_id"`
	AliasTargetEvaluateTargetHealth types.Bool   `tfsdk:"alias_target_evaluate_target_health"`
	Failover                        types.String `tfsdk:"failover"`
	GeolocationContinentCode        types.String `tfsdk:"geolocation_continent_code"`
	GeolocationCountryCode          types.String `tfsdk:"geolocation_country_code"`
	GeolocationSubdivisionCode      types.String `tfsdk:"geolocation_subdivision_code"`
	HealthCheckId                   types.String `tfsdk:"health_check_id"`
	Region                          types.String `tfsdk:"region"`
	SetIdentifier                   types.String `tfsdk:"set_identifier"`
	Type                            types.String `tfsdk:"type"`
	Weight                          types.Int64  `tfsdk:"weight"`
}

var RecordAAwsRte53RecordInfoAttrTypes = map[string]attr.Type{
	"alias_target_dns_name":               types.StringType,
	"alias_target_hosted_zone_id":         types.StringType,
	"alias_target_evaluate_target_health": types.BoolType,
	"failover":                            types.StringType,
	"geolocation_continent_code":          types.StringType,
	"geolocation_country_code":            types.StringType,
	"geolocation_subdivision_code":        types.StringType,
	"health_check_id":                     types.StringType,
	"region":                              types.StringType,
	"set_identifier":                      types.StringType,
	"type":                                types.StringType,
	"weight":                              types.Int64Type,
}

var RecordAAwsRte53RecordInfoResourceSchemaAttributes = map[string]schema.Attribute{
	"alias_target_dns_name": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "DNS name of the alias target.",
	},
	"alias_target_hosted_zone_id": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Hosted zone ID of the alias target.",
	},
	"alias_target_evaluate_target_health": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Indicates if Amazon Route 53 evaluates the health of the alias target.",
	},
	"failover": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Indicates whether this is the primary or secondary resource record for Amazon Route 53 failover routing.",
	},
	"geolocation_continent_code": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Continent code for Amazon Route 53 geolocation routing.",
	},
	"geolocation_country_code": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Country code for Amazon Route 53 geolocation routing.",
	},
	"geolocation_subdivision_code": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Subdivision code for Amazon Route 53 geolocation routing.",
	},
	"health_check_id": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "ID of the health check that Amazon Route 53 performs for this resource record.",
	},
	"region": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Amazon EC2 region where this resource record resides for latency routing.",
	},
	"set_identifier": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "An identifier that differentiates records with the same DNS name and type for weighted, latency, geolocation, and failover routing.",
	},
	"type": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Type of Amazon Route 53 resource record.",
	},
	"weight": schema.Int64Attribute{
		Computed:            true,
		MarkdownDescription: "Value that determines the portion of traffic for this record in weighted routing. The range is from 0 to 255.",
	},
}

func ExpandRecordAAwsRte53RecordInfo(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.RecordAAwsRte53RecordInfo {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m RecordAAwsRte53RecordInfoModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *RecordAAwsRte53RecordInfoModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.RecordAAwsRte53RecordInfo {
	if m == nil {
		return nil
	}
	to := &dns.RecordAAwsRte53RecordInfo{}
	return to
}

func FlattenRecordAAwsRte53RecordInfo(ctx context.Context, from *dns.RecordAAwsRte53RecordInfo, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RecordAAwsRte53RecordInfoAttrTypes)
	}
	m := RecordAAwsRte53RecordInfoModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RecordAAwsRte53RecordInfoAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RecordAAwsRte53RecordInfoModel) Flatten(ctx context.Context, from *dns.RecordAAwsRte53RecordInfo, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RecordAAwsRte53RecordInfoModel{}
	}
	m.AliasTargetDnsName = flex.FlattenStringPointer(from.AliasTargetDnsName)
	m.AliasTargetHostedZoneId = flex.FlattenStringPointer(from.AliasTargetHostedZoneId)
	m.AliasTargetEvaluateTargetHealth = types.BoolPointerValue(from.AliasTargetEvaluateTargetHealth)
	m.Failover = flex.FlattenStringPointer(from.Failover)
	m.GeolocationContinentCode = flex.FlattenStringPointer(from.GeolocationContinentCode)
	m.GeolocationCountryCode = flex.FlattenStringPointer(from.GeolocationCountryCode)
	m.GeolocationSubdivisionCode = flex.FlattenStringPointer(from.GeolocationSubdivisionCode)
	m.HealthCheckId = flex.FlattenStringPointer(from.HealthCheckId)
	m.Region = flex.FlattenStringPointer(from.Region)
	m.SetIdentifier = flex.FlattenStringPointer(from.SetIdentifier)
	m.Type = flex.FlattenStringPointer(from.Type)
	m.Weight = flex.FlattenInt64Pointer(from.Weight)
}

func (m *RecordAAwsRte53RecordInfoModel) PutExpand(to *dns.RecordAAwsRte53RecordInfo) *dns.RecordAAwsRte53RecordInfo {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RecordAAwsRte53RecordInfoResourceSchemaAttributes {
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
