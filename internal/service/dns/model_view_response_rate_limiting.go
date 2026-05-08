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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ViewResponseRateLimitingModel struct {
	EnableRrl          types.Bool  `tfsdk:"enable_rrl"`
	LogOnly            types.Bool  `tfsdk:"log_only"`
	ResponsesPerSecond types.Int64 `tfsdk:"responses_per_second"`
	Window             types.Int64 `tfsdk:"window"`
	Slip               types.Int64 `tfsdk:"slip"`
}

var ViewResponseRateLimitingAttrTypes = map[string]attr.Type{
	"enable_rrl":           types.BoolType,
	"log_only":             types.BoolType,
	"responses_per_second": types.Int64Type,
	"window":               types.Int64Type,
	"slip":                 types.Int64Type,
}

var ViewResponseRateLimitingResourceSchemaAttributes = map[string]schema.Attribute{
	"enable_rrl": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if the response rate limiting is enabled or not.",
	},
	"log_only": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if logging for response rate limiting without dropping any requests is enabled or not.",
	},
	"responses_per_second": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of responses per client per second.",
	},
	"window": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The time interval in seconds over which responses are tracked.",
	},
	"slip": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The response rate limiting slip. Note that if slip is not equal to 0 every n-th rate-limited UDP request is sent a truncated response instead of being dropped.",
	},
}

func ExpandViewResponseRateLimiting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ViewResponseRateLimiting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ViewResponseRateLimitingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ViewResponseRateLimitingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ViewResponseRateLimiting {
	if m == nil {
		return nil
	}
	to := &dns.ViewResponseRateLimiting{
		EnableRrl:          flex.ExpandBoolPointer(m.EnableRrl),
		LogOnly:            flex.ExpandBoolPointer(m.LogOnly),
		ResponsesPerSecond: flex.ExpandInt64Pointer(m.ResponsesPerSecond),
		Window:             flex.ExpandInt64Pointer(m.Window),
		Slip:               flex.ExpandInt64Pointer(m.Slip),
	}
	return to
}

func FlattenViewResponseRateLimiting(ctx context.Context, from *dns.ViewResponseRateLimiting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ViewResponseRateLimitingAttrTypes)
	}
	m := ViewResponseRateLimitingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ViewResponseRateLimitingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ViewResponseRateLimitingModel) Flatten(ctx context.Context, from *dns.ViewResponseRateLimiting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ViewResponseRateLimitingModel{}
	}
	m.EnableRrl = types.BoolPointerValue(from.EnableRrl)
	m.LogOnly = types.BoolPointerValue(from.LogOnly)
	m.ResponsesPerSecond = flex.FlattenInt64Pointer(from.ResponsesPerSecond)
	m.Window = flex.FlattenInt64Pointer(from.Window)
	m.Slip = flex.FlattenInt64Pointer(from.Slip)
}

func (m *ViewResponseRateLimitingModel) PutExpand(to *dns.ViewResponseRateLimiting) *dns.ViewResponseRateLimiting {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ViewResponseRateLimitingResourceSchemaAttributes {
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
