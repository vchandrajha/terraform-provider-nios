package ipam

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type VlanModel struct {
	Ref         types.String `tfsdk:"ref"`
	AssignedTo  types.List   `tfsdk:"assigned_to"`
	Comment     types.String `tfsdk:"comment"`
	Contact     types.String `tfsdk:"contact"`
	Department  types.String `tfsdk:"department"`
	Description types.String `tfsdk:"description"`
	ExtAttrs    types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll types.Map    `tfsdk:"extattrs_all"`
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Parent      types.String `tfsdk:"parent"`
	Reserved    types.Bool   `tfsdk:"reserved"`
	Status      types.String `tfsdk:"status"`
}

var VlanAttrTypes = map[string]attr.Type{
	"ref":          types.StringType,
	"assigned_to":  types.ListType{ElemType: types.StringType},
	"comment":      types.StringType,
	"contact":      types.StringType,
	"department":   types.StringType,
	"description":  types.StringType,
	"extattrs":     types.MapType{ElemType: types.StringType},
	"extattrs_all": types.MapType{ElemType: types.StringType},
	"id":           types.Int64Type,
	"name":         types.StringType,
	"parent":       types.StringType,
	"reserved":     types.BoolType,
	"status":       types.StringType,
}

var VlanResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"assigned_to": schema.ListAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "List of objects VLAN is assigned to.",
	},
	"comment": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthBetween(0, 256),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "A descriptive comment for this VLAN.",
	},
	"contact": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "Contact information for person/team managing or using VLAN.",
	},
	"department": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "Department where VLAN is used.",
	},
	"description": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "Description for the VLAN object, may be potentially used for longer VLAN names.",
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
		Computed:            true,
		MarkdownDescription: "Extensible attributes associated with the object , including default and internal attributes.",
		ElementType:         types.StringType,
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
	},
	"id": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "VLAN ID value.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Name of the VLAN.",
	},
	"parent": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The VLAN View or VLAN Range to which this VLAN belongs.",
	},
	"reserved": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "When set VLAN can only be assigned to IPAM object manually.",
	},
	"status": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Status of VLAN object. Can be Assigned, Unassigned, Reserved.",
	},
}

func (m *VlanModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.Vlan {
	if m == nil {
		return nil
	}
	to := &ipam.Vlan{
		Comment:     flex.ExpandStringPointer(m.Comment),
		Contact:     flex.ExpandStringPointer(m.Contact),
		Department:  flex.ExpandStringPointer(m.Department),
		Description: flex.ExpandStringPointer(m.Description),
		ExtAttrs:    ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Id:          ExpandVlanId(m.Id),
		Name:        flex.ExpandStringPointer(m.Name),
		Parent:      ExpandVlanParent(m.Parent),
		Reserved:    flex.ExpandBoolPointer(m.Reserved),
	}
	return to
}

func FlattenVlan(ctx context.Context, from *ipam.Vlan, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(VlanAttrTypes)
	}
	m := VlanModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, VlanAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *VlanModel) Flatten(ctx context.Context, from *ipam.Vlan, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = VlanModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AssignedTo = flex.FlattenFrameworkListString(ctx, from.AssignedTo, diags)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Contact = flex.FlattenStringPointer(from.Contact)
	m.Department = flex.FlattenStringPointer(from.Department)
	m.Description = flex.FlattenStringPointer(from.Description)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Id = FlattenVlanId(from.Id)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Parent = FlattenVlanParent(from.Parent)
	m.Reserved = types.BoolPointerValue(from.Reserved)
	m.Status = flex.FlattenStringPointer(from.Status)
}

func ExpandVlanParent(str types.String) *ipam.VlanParent {
	if str.IsNull() {
		return &ipam.VlanParent{}
	}
	var m ipam.VlanParent
	m.String = flex.ExpandStringPointer(str)

	return &m
}

func FlattenVlanParent(from *ipam.VlanParent) types.String {
	if from.VlanParentOneOf == nil {
		return types.StringNull()
	}
	m := flex.FlattenStringPointer(from.VlanParentOneOf.Ref)
	return m
}

func ExpandVlanId(val types.Int64) *ipam.VlanId {
	if val.IsNull() {
		return &ipam.VlanId{}
	}
	var m ipam.VlanId
	m.Int64 = flex.ExpandInt64Pointer(val)

	return &m
}

func FlattenVlanId(from *ipam.VlanId) types.Int64 {
	if from.Int64 == nil {
		return types.Int64Null()
	}
	m := flex.FlattenInt64Pointer(from.Int64)
	return m
}

func (m *VlanModel) PutExpand(to *ipam.Vlan) *ipam.Vlan {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range VlanResourceSchemaAttributes {
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
