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

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneRpMemberSoaSerialsModel struct {
	GridPrimary     types.String `tfsdk:"grid_primary"`
	MsServerPrimary types.String `tfsdk:"ms_server_primary"`
	Serial          types.Int64  `tfsdk:"serial"`
}

var ZoneRpMemberSoaSerialsAttrTypes = map[string]attr.Type{
	"grid_primary":      types.StringType,
	"ms_server_primary": types.StringType,
	"serial":            types.Int64Type,
}

var ZoneRpMemberSoaSerialsResourceSchemaAttributes = map[string]schema.Attribute{
	"grid_primary": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The grid primary server for the zone. Only one of \"grid_primary\" or \"ms_server_primary\" will be set when the object is retrieved from the server.",
	},
	"ms_server_primary": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The primary MS server for the zone. Only one of \"grid_primary\" or \"ms_server_primary\" will be set when the object is retrieved from the server.",
	},
	"serial": schema.Int64Attribute{
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The SOA serial number.",
	},
}

func ExpandZoneRpMemberSoaSerials(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneRpMemberSoaSerials {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneRpMemberSoaSerialsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneRpMemberSoaSerialsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneRpMemberSoaSerials {
	if m == nil {
		return nil
	}
	to := &dns.ZoneRpMemberSoaSerials{}
	return to
}

func FlattenZoneRpMemberSoaSerials(ctx context.Context, from *dns.ZoneRpMemberSoaSerials, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneRpMemberSoaSerialsAttrTypes)
	}
	m := ZoneRpMemberSoaSerialsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneRpMemberSoaSerialsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneRpMemberSoaSerialsModel) Flatten(ctx context.Context, from *dns.ZoneRpMemberSoaSerials, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneRpMemberSoaSerialsModel{}
	}
	m.GridPrimary = flex.FlattenStringPointer(from.GridPrimary)
	m.MsServerPrimary = flex.FlattenStringPointer(from.MsServerPrimary)
	m.Serial = flex.FlattenInt64Pointer(from.Serial)
}

func (m *ZoneRpMemberSoaSerialsModel) PutExpand(to *dns.ZoneRpMemberSoaSerials) *dns.ZoneRpMemberSoaSerials {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneRpMemberSoaSerialsResourceSchemaAttributes {
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
