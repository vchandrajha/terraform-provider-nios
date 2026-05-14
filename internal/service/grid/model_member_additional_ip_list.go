package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberAdditionalIpListModel struct {
	Anycast            types.Bool   `tfsdk:"anycast"`
	Ipv4NetworkSetting types.Object `tfsdk:"ipv4_network_setting"`
	Ipv6NetworkSetting types.Object `tfsdk:"ipv6_network_setting"`
	Comment            types.String `tfsdk:"comment"`
	EnableBgp          types.Bool   `tfsdk:"enable_bgp"`
	EnableOspf         types.Bool   `tfsdk:"enable_ospf"`
	Interface          types.String `tfsdk:"interface"`
}

var MemberAdditionalIpListAttrTypes = map[string]attr.Type{
	"anycast":              types.BoolType,
	"ipv4_network_setting": types.ObjectType{AttrTypes: MemberadditionaliplistIpv4NetworkSettingAttrTypes},
	"ipv6_network_setting": types.ObjectType{AttrTypes: MemberadditionaliplistIpv6NetworkSettingAttrTypes},
	"comment":              types.StringType,
	"enable_bgp":           types.BoolType,
	"enable_ospf":          types.BoolType,
	"interface":            types.StringType,
}

var MemberAdditionalIpListResourceSchemaAttributes = map[string]schema.Attribute{
	"anycast": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if anycast for the Interface object is enabled or not.",
	},
	"ipv4_network_setting": schema.SingleNestedAttribute{
		Attributes:          MemberadditionaliplistIpv4NetworkSettingResourceSchemaAttributes,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The IPv4 network settings of the Grid Member.",
	},
	"ipv6_network_setting": schema.SingleNestedAttribute{
		Attributes:          MemberadditionaliplistIpv6NetworkSettingResourceSchemaAttributes,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The IPv6 network settings of the Grid Member.",
	},
	"comment": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "A descriptive comment of this structure.",
	},
	"enable_bgp": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the BGP advertisement setting is enabled for this interface or not.",
	},
	"enable_ospf": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the OSPF advertisement setting is enabled for this interface or not.",
	},
	"interface": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("LOOPBACK"),
		Validators: []validator.String{
			stringvalidator.OneOf("LAN2", "LAN_HA", "LOOPBACK", "MGMT"),
		},
		MarkdownDescription: "The interface type for the Interface object.",
	},
}

func ExpandMemberAdditionalIpList(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberAdditionalIpList {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberAdditionalIpListModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberAdditionalIpListModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberAdditionalIpList {
	if m == nil {
		return nil
	}
	to := &grid.MemberAdditionalIpList{
		Anycast:            flex.ExpandBoolPointer(m.Anycast),
		Ipv4NetworkSetting: ExpandMemberadditionaliplistIpv4NetworkSetting(ctx, m.Ipv4NetworkSetting, diags),
		Ipv6NetworkSetting: ExpandMemberadditionaliplistIpv6NetworkSetting(ctx, m.Ipv6NetworkSetting, diags),
		Comment:            flex.ExpandStringPointer(m.Comment),
		EnableBgp:          flex.ExpandBoolPointer(m.EnableBgp),
		EnableOspf:         flex.ExpandBoolPointer(m.EnableOspf),
		Interface:          flex.ExpandStringPointer(m.Interface),
	}
	return to
}

func FlattenMemberAdditionalIpList(ctx context.Context, from *grid.MemberAdditionalIpList, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberAdditionalIpListAttrTypes)
	}
	m := MemberAdditionalIpListModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberAdditionalIpListAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberAdditionalIpListModel) Flatten(ctx context.Context, from *grid.MemberAdditionalIpList, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberAdditionalIpListModel{}
	}
	m.Anycast = types.BoolPointerValue(from.Anycast)
	m.Ipv4NetworkSetting = FlattenMemberadditionaliplistIpv4NetworkSetting(ctx, from.Ipv4NetworkSetting, diags)
	m.Ipv6NetworkSetting = FlattenMemberadditionaliplistIpv6NetworkSetting(ctx, from.Ipv6NetworkSetting, diags)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.EnableBgp = types.BoolPointerValue(from.EnableBgp)
	m.EnableOspf = types.BoolPointerValue(from.EnableOspf)
	m.Interface = flex.FlattenStringPointer(from.Interface)
}

func (m *MemberAdditionalIpListModel) PutExpand(to *grid.MemberAdditionalIpList) *grid.MemberAdditionalIpList {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberAdditionalIpListResourceSchemaAttributes {
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
