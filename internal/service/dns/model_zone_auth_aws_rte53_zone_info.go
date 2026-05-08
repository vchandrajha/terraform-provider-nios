package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneAuthAwsRte53ZoneInfoModel struct {
	AssociatedVpcs  types.List   `tfsdk:"associated_vpcs"`
	CallerReference types.String `tfsdk:"callerreference"`
	DelegationSetId types.String `tfsdk:"delegation_set_id"`
	HostedZoneId    types.String `tfsdk:"hosted_zone_id"`
	NameServers     types.List   `tfsdk:"name_servers"`
	RecordSetCount  types.Int64  `tfsdk:"record_set_count"`
	Type            types.String `tfsdk:"type"`
}

var ZoneAuthAwsRte53ZoneInfoAttrTypes = map[string]attr.Type{
	"associated_vpcs":   types.ListType{ElemType: types.StringType},
	"callerreference":   types.StringType,
	"delegation_set_id": types.StringType,
	"hosted_zone_id":    types.StringType,
	"name_servers":      types.ListType{ElemType: types.StringType},
	"record_set_count":  types.Int64Type,
	"type":              types.StringType,
}

var ZoneAuthAwsRte53ZoneInfoResourceSchemaAttributes = map[string]schema.Attribute{
	"associated_vpcs": schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "List of AWS VPC strings that are associated with this zone.",
	},
	"callerreference": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "User specified caller reference when zone was created.",
	},
	"delegation_set_id": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "ID of delegation set associated with this zone.",
	},
	"hosted_zone_id": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "AWS route 53 assigned ID for this zone.",
	},
	"name_servers": schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "List of AWS name servers that are authoritative for this domain name.",
	},
	"record_set_count": schema.Int64Attribute{
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Number of resource record sets in the hosted zone.",
	},
	"type": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Indicates whether private or public zone.",
	},
}

func ExpandZoneAuthAwsRte53ZoneInfo(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthAwsRte53ZoneInfo {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthAwsRte53ZoneInfoModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthAwsRte53ZoneInfoModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthAwsRte53ZoneInfo {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthAwsRte53ZoneInfo{}
	return to
}

func FlattenZoneAuthAwsRte53ZoneInfo(ctx context.Context, from *dns.ZoneAuthAwsRte53ZoneInfo, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthAwsRte53ZoneInfoAttrTypes)
	}
	m := ZoneAuthAwsRte53ZoneInfoModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthAwsRte53ZoneInfoAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthAwsRte53ZoneInfoModel) Flatten(ctx context.Context, from *dns.ZoneAuthAwsRte53ZoneInfo, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthAwsRte53ZoneInfoModel{}
	}
	m.AssociatedVpcs = flex.FlattenFrameworkListString(ctx, from.AssociatedVpcs, diags)
	m.CallerReference = flex.FlattenStringPointer(from.CallerReference)
	m.DelegationSetId = flex.FlattenStringPointer(from.DelegationSetId)
	m.HostedZoneId = flex.FlattenStringPointer(from.HostedZoneId)
	m.NameServers = flex.FlattenFrameworkListString(ctx, from.NameServers, diags)
	m.RecordSetCount = flex.FlattenInt64Pointer(from.RecordSetCount)
	m.Type = flex.FlattenStringPointer(from.Type)
}

func (m *ZoneAuthAwsRte53ZoneInfoModel) PutExpand(to *dns.ZoneAuthAwsRte53ZoneInfo) *dns.ZoneAuthAwsRte53ZoneInfo {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthAwsRte53ZoneInfoResourceSchemaAttributes {
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
