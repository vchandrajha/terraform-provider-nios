package microsoft

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"

	"github.com/infobloxopen/infoblox-nios-go-client/microsoft"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type MssuperscopeModel struct {
	Ref                   types.String `tfsdk:"ref"`
	Comment               types.String `tfsdk:"comment"`
	DhcpUtilization       types.Int64  `tfsdk:"dhcp_utilization"`
	DhcpUtilizationStatus types.String `tfsdk:"dhcp_utilization_status"`
	Disable               types.Bool   `tfsdk:"disable"`
	DynamicHosts          types.Int64  `tfsdk:"dynamic_hosts"`
	ExtAttrs              types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll           types.Map    `tfsdk:"extattrs_all"`
	HighWaterMark         types.Int64  `tfsdk:"high_water_mark"`
	HighWaterMarkReset    types.Int64  `tfsdk:"high_water_mark_reset"`
	LowWaterMark          types.Int64  `tfsdk:"low_water_mark"`
	LowWaterMarkReset     types.Int64  `tfsdk:"low_water_mark_reset"`
	Name                  types.String `tfsdk:"name"`
	NetworkView           types.String `tfsdk:"network_view"`
	Ranges                types.List   `tfsdk:"ranges"`
	StaticHosts           types.Int64  `tfsdk:"static_hosts"`
	TotalHosts            types.Int64  `tfsdk:"total_hosts"`
}

var MssuperscopeAttrTypes = map[string]attr.Type{
	"ref":                     types.StringType,
	"comment":                 types.StringType,
	"dhcp_utilization":        types.Int64Type,
	"dhcp_utilization_status": types.StringType,
	"disable":                 types.BoolType,
	"dynamic_hosts":           types.Int64Type,
	"extattrs":                types.MapType{ElemType: types.StringType},
	"extattrs_all":            types.MapType{ElemType: types.StringType},
	"high_water_mark":         types.Int64Type,
	"high_water_mark_reset":   types.Int64Type,
	"low_water_mark":          types.Int64Type,
	"low_water_mark_reset":    types.Int64Type,
	"name":                    types.StringType,
	"network_view":            types.StringType,
	"ranges":                  types.ListType{ElemType: types.StringType},
	"static_hosts":            types.Int64Type,
	"total_hosts":             types.Int64Type,
}

var MssuperscopeResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The superscope descriptive comment.",
	},
	"dhcp_utilization": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The percentage of the total DHCP usage of the ranges in the superscope.",
	},
	"dhcp_utilization_status": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Utilization level of the DHCP range objects that belong to the superscope.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines whether the superscope is disabled.",
	},
	"dynamic_hosts": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The total number of DHCP leases issued for the DHCP range objects that belong to the superscope.",
	},
	"extattrs": schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Extensible attributes associated with the object.",
	},
	"extattrs_all": schema.MapAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Extensible attributes associated with the object, including default and internal attributes.",
		ElementType:         types.StringType,
	},
	"high_water_mark": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The percentage value for DHCP range usage after which an alarm will be active.",
	},
	"high_water_mark_reset": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The percentage value for DHCP range usage after which an alarm will be reset.",
	},
	"low_water_mark": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The percentage value for DHCP range usage below which an alarm will be active.",
	},
	"low_water_mark_reset": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The percentage value for DHCP range usage below which an alarm will be reset.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of the Microsoft DHCP superscope.",
	},
	"network_view": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString("default"),
		MarkdownDescription: "The name of the network view in which the superscope resides.",
	},
	"ranges": schema.ListAttribute{
		ElementType: types.StringType,
		Required:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The list of DHCP ranges that are associated with the superscope.",
	},
	"static_hosts": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of static DHCP addresses configured in DHCP range objects that belong to the superscope.",
	},
	"total_hosts": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The total number of DHCP addresses configured in DHCP range objects that belong to the superscope.",
	},
}

func ExpandMssuperscope(ctx context.Context, o types.Object, diags *diag.Diagnostics) *microsoft.Mssuperscope {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MssuperscopeModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MssuperscopeModel) Expand(ctx context.Context, diags *diag.Diagnostics) *microsoft.Mssuperscope {
	if m == nil {
		return nil
	}
	to := &microsoft.Mssuperscope{
		Comment:     flex.ExpandStringPointer(m.Comment),
		Disable:     flex.ExpandBoolPointer(m.Disable),
		ExtAttrs:    ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Name:        flex.ExpandStringPointer(m.Name),
		NetworkView: flex.ExpandStringPointer(m.NetworkView),
		Ranges:      flex.ExpandFrameworkListString(ctx, m.Ranges, diags),
	}
	return to
}

func FlattenMssuperscope(ctx context.Context, from *microsoft.Mssuperscope, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MssuperscopeAttrTypes)
	}
	m := MssuperscopeModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, MssuperscopeAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MssuperscopeModel) Flatten(ctx context.Context, from *microsoft.Mssuperscope, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MssuperscopeModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.DhcpUtilization = flex.FlattenInt64Pointer(from.DhcpUtilization)
	m.DhcpUtilizationStatus = flex.FlattenStringPointer(from.DhcpUtilizationStatus)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.DynamicHosts = flex.FlattenInt64Pointer(from.DynamicHosts)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.HighWaterMark = flex.FlattenInt64Pointer(from.HighWaterMark)
	m.HighWaterMarkReset = flex.FlattenInt64Pointer(from.HighWaterMarkReset)
	m.LowWaterMark = flex.FlattenInt64Pointer(from.LowWaterMark)
	m.LowWaterMarkReset = flex.FlattenInt64Pointer(from.LowWaterMarkReset)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.NetworkView = flex.FlattenStringPointer(from.NetworkView)
	m.Ranges = flex.FlattenFrameworkListString(ctx, from.Ranges, diags)
	m.StaticHosts = flex.FlattenInt64Pointer(from.StaticHosts)
	m.TotalHosts = flex.FlattenInt64Pointer(from.TotalHosts)
}

func (m *MssuperscopeModel) PutExpand(to *microsoft.Mssuperscope) *microsoft.Mssuperscope {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MssuperscopeResourceSchemaAttributes {
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
