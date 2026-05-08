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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneForwardForwardingServersModel struct {
	Name                  types.String `tfsdk:"name"`
	ForwardersOnly        types.Bool   `tfsdk:"forwarders_only"`
	ForwardTo             types.List   `tfsdk:"forward_to"`
	UseOverrideForwarders types.Bool   `tfsdk:"use_override_forwarders"`
}

var ZoneForwardForwardingServersAttrTypes = map[string]attr.Type{
	"name":                    types.StringType,
	"forwarders_only":         types.BoolType,
	"forward_to":              types.ListType{ElemType: types.ObjectType{AttrTypes: ZoneforwardforwardingserversForwardToAttrTypes}},
	"use_override_forwarders": types.BoolType,
}

var ZoneForwardForwardingServersResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name of this Grid member in FQDN format.",
	},
	"forwarders_only": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the appliance sends queries to forwarders only, and not to other internal or Internet root servers.",
	},
	"forward_to": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: ZoneforwardforwardingserversForwardToResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The information for the remote name server to which you want the Infoblox appliance to forward queries for a specified domain name.",
	},
	"use_override_forwarders": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Use flag for: forward_to",
	},
}

func ExpandZoneForwardForwardingServers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneForwardForwardingServers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneForwardForwardingServersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneForwardForwardingServersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneForwardForwardingServers {
	if m == nil {
		return nil
	}
	to := &dns.ZoneForwardForwardingServers{
		Name:                  flex.ExpandStringPointer(m.Name),
		ForwardersOnly:        flex.ExpandBoolPointer(m.ForwardersOnly),
		ForwardTo:             flex.ExpandFrameworkListNestedBlock(ctx, m.ForwardTo, diags, ExpandZoneforwardforwardingserversForwardTo),
		UseOverrideForwarders: flex.ExpandBoolPointer(m.UseOverrideForwarders),
	}
	return to
}

func FlattenZoneForwardForwardingServers(ctx context.Context, from *dns.ZoneForwardForwardingServers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneForwardForwardingServersAttrTypes)
	}
	m := ZoneForwardForwardingServersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneForwardForwardingServersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneForwardForwardingServersModel) Flatten(ctx context.Context, from *dns.ZoneForwardForwardingServers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneForwardForwardingServersModel{}
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.ForwardersOnly = types.BoolPointerValue(from.ForwardersOnly)
	m.ForwardTo = flex.FlattenFrameworkListNestedBlock(ctx, from.ForwardTo, ZoneforwardforwardingserversForwardToAttrTypes, diags, FlattenZoneforwardforwardingserversForwardTo)
	m.UseOverrideForwarders = types.BoolPointerValue(from.UseOverrideForwarders)
}

func (m *ZoneForwardForwardingServersModel) PutExpand(to *dns.ZoneForwardForwardingServers) *dns.ZoneForwardForwardingServers {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneForwardForwardingServersResourceSchemaAttributes {
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
