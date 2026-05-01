package ipam

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

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type NetworkviewDdnsZonePrimariesModel struct {
	ZoneMatch      types.String `tfsdk:"zone_match"`
	DnsGridZone    types.Object `tfsdk:"dns_grid_zone"`
	DnsGridPrimary types.String `tfsdk:"dns_grid_primary"`
	DnsExtZone     types.String `tfsdk:"dns_ext_zone"`
	DnsExtPrimary  types.String `tfsdk:"dns_ext_primary"`
}

var NetworkviewDdnsZonePrimariesAttrTypes = map[string]attr.Type{
	"zone_match":       types.StringType,
	"dns_grid_zone":    types.ObjectType{AttrTypes: NetworkviewDdnsZonePrimariesDnsGridZoneAttrTypes},
	"dns_grid_primary": types.StringType,
	"dns_ext_zone":     types.StringType,
	"dns_ext_primary":  types.StringType,
}

var NetworkviewDdnsZonePrimariesResourceSchemaAttributes = map[string]schema.Attribute{
	"zone_match": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ANY_EXTERNAL", "ANY_GRID", "EXTERNAL", "GRID"),
		},
		MarkdownDescription: "Indicate matching type.",
	},
	"dns_grid_zone": schema.SingleNestedAttribute{
		Attributes:          NetworkviewDdnsZonePrimariesDnsGridZoneResourceSchemaAttributes,
		Optional:            true,
		MarkdownDescription: "The ref of a DNS zone.",
	},
	"dns_grid_primary": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The name of a Grid member.",
	},
	"dns_ext_zone": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.IsValidFQDN(),
		},
		MarkdownDescription: "The name of external zone in FQDN format.",
	},
	"dns_ext_primary": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The IP address of the External server. Valid when zone_match is \"EXTERNAL\" or \"ANY_EXTERNAL\".",
	},
}

func ExpandNetworkviewDdnsZonePrimaries(ctx context.Context, o types.Object, diags *diag.Diagnostics) *ipam.NetworkviewDdnsZonePrimaries {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m NetworkviewDdnsZonePrimariesModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *NetworkviewDdnsZonePrimariesModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.NetworkviewDdnsZonePrimaries {
	if m == nil {
		return nil
	}
	to := &ipam.NetworkviewDdnsZonePrimaries{
		ZoneMatch:      flex.ExpandStringPointer(m.ZoneMatch),
		DnsGridZone:    ExpandNetworkviewDdnsZonePrimariesDnsGridZone(ctx, m.DnsGridZone, diags),
		DnsGridPrimary: flex.ExpandStringPointer(m.DnsGridPrimary),
		DnsExtZone:     flex.ExpandStringPointer(m.DnsExtZone),
		DnsExtPrimary:  flex.ExpandStringPointer(m.DnsExtPrimary),
	}
	return to
}

func FlattenNetworkviewDdnsZonePrimaries(ctx context.Context, from *ipam.NetworkviewDdnsZonePrimaries, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NetworkviewDdnsZonePrimariesAttrTypes)
	}
	m := NetworkviewDdnsZonePrimariesModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, NetworkviewDdnsZonePrimariesAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NetworkviewDdnsZonePrimariesModel) Flatten(ctx context.Context, from *ipam.NetworkviewDdnsZonePrimaries, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NetworkviewDdnsZonePrimariesModel{}
	}
	m.ZoneMatch = flex.FlattenStringPointer(from.ZoneMatch)
	m.DnsGridZone = FlattenNetworkviewDdnsZonePrimariesDnsGridZone(ctx, from.DnsGridZone, diags)
	m.DnsGridPrimary = flex.FlattenStringPointer(from.DnsGridPrimary)
	m.DnsExtZone = flex.FlattenStringPointer(from.DnsExtZone)
	m.DnsExtPrimary = flex.FlattenStringPointer(from.DnsExtPrimary)
}

func (m *NetworkviewDdnsZonePrimariesModel) PutExpand(to *ipam.NetworkviewDdnsZonePrimaries) *ipam.NetworkviewDdnsZonePrimaries {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NetworkviewDdnsZonePrimariesResourceSchemaAttributes {
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
