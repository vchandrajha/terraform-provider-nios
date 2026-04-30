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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneAuthDnssecKeysModel struct {
	Tag           types.Int64  `tfsdk:"tag"`
	Status        types.String `tfsdk:"status"`
	NextEventDate types.Int64  `tfsdk:"next_event_date"`
	Type          types.String `tfsdk:"type"`
	Algorithm     types.String `tfsdk:"algorithm"`
	PublicKey     types.String `tfsdk:"public_key"`
}

var ZoneAuthDnssecKeysAttrTypes = map[string]attr.Type{
	"tag":             types.Int64Type,
	"status":          types.StringType,
	"next_event_date": types.Int64Type,
	"type":            types.StringType,
	"algorithm":       types.StringType,
	"public_key":      types.StringType,
}

var ZoneAuthDnssecKeysResourceSchemaAttributes = map[string]schema.Attribute{
	"tag": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The tag of the key for the zone.",
	},
	"status": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The status of the key for the zone.",
	},
	"next_event_date": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The next event date for the key, the rollover date for an active key or the removal date for an already rolled one.",
	},
	"type": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The key type.",
	},
	"algorithm": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The public-key encryption algorithm. Values 1, 3 and 6 are deprecated from NIOS 9.0.",
	},
	"public_key": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The Base-64 encoding of the public key.",
	},
}

func ExpandZoneAuthDnssecKeys(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthDnssecKeys {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthDnssecKeysModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthDnssecKeysModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthDnssecKeys {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthDnssecKeys{
		Tag: flex.ExpandInt64Pointer(m.Tag),
	}
	return to
}

func FlattenZoneAuthDnssecKeys(ctx context.Context, from *dns.ZoneAuthDnssecKeys, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthDnssecKeysAttrTypes)
	}
	m := ZoneAuthDnssecKeysModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthDnssecKeysAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthDnssecKeysModel) Flatten(ctx context.Context, from *dns.ZoneAuthDnssecKeys, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthDnssecKeysModel{}
	}
	m.Tag = flex.FlattenInt64Pointer(from.Tag)
	m.Status = flex.FlattenStringPointer(from.Status)
	m.NextEventDate = flex.FlattenInt64Pointer(from.NextEventDate)
	m.Type = flex.FlattenStringPointer(from.Type)
	m.Algorithm = flex.FlattenStringPointer(from.Algorithm)
	m.PublicKey = flex.FlattenStringPointer(from.PublicKey)
}

func (m *ZoneAuthDnssecKeysModel) PutExpand(to *dns.ZoneAuthDnssecKeys) *dns.ZoneAuthDnssecKeys {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthDnssecKeysResourceSchemaAttributes {
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
