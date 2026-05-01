package grid

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type UpgradegroupModel struct {
	Ref                        types.String `tfsdk:"ref"`
	Comment                    types.String `tfsdk:"comment"`
	DistributionDependentGroup types.String `tfsdk:"distribution_dependent_group"`
	DistributionPolicy         types.String `tfsdk:"distribution_policy"`
	DistributionTime           types.String `tfsdk:"distribution_time"`
	Members                    types.List   `tfsdk:"members"`
	Name                       types.String `tfsdk:"name"`
	TimeZone                   types.String `tfsdk:"time_zone"`
	UpgradeDependentGroup      types.String `tfsdk:"upgrade_dependent_group"`
	UpgradePolicy              types.String `tfsdk:"upgrade_policy"`
	UpgradeTime                types.String `tfsdk:"upgrade_time"`
}

var UpgradegroupAttrTypes = map[string]attr.Type{
	"ref":                          types.StringType,
	"comment":                      types.StringType,
	"distribution_dependent_group": types.StringType,
	"distribution_policy":          types.StringType,
	"distribution_time":            types.StringType,
	"members":                      types.ListType{ElemType: types.ObjectType{AttrTypes: UpgradegroupMembersAttrTypes}},
	"name":                         types.StringType,
	"time_zone":                    types.StringType,
	"upgrade_dependent_group":      types.StringType,
	"upgrade_policy":               types.StringType,
	"upgrade_time":                 types.StringType,
}

var UpgradegroupResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The upgrade group descriptive comment.",
	},
	"distribution_dependent_group": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The distribution dependent group name.",
	},
	"distribution_policy": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("SIMULTANEOUSLY"),
		Validators: []validator.String{
			stringvalidator.OneOf("SIMULTANEOUSLY", "SEQUENTIALLY"),
		},
		MarkdownDescription: "The distribution scheduling policy.",
	},
	"distribution_time": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.ValidateTimeFormat(),
		},
		MarkdownDescription: "The time of the next scheduled distribution.",
	},
	"members": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: UpgradegroupMembersResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The upgrade group members.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The upgrade group name.",
	},
	"time_zone": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The time zone for scheduling operations.",
	},
	"upgrade_dependent_group": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The upgrade dependent group name.",
	},
	"upgrade_policy": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("SEQUENTIALLY"),
		Validators: []validator.String{
			stringvalidator.OneOf("SIMULTANEOUSLY", "SEQUENTIALLY"),
		},
		MarkdownDescription: "The upgrade scheduling policy.",
	},
	"upgrade_time": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			customvalidator.ValidateTimeFormat(),
		},
		MarkdownDescription: "The time of the next scheduled upgrade.",
	},
}

func (m *UpgradegroupModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.Upgradegroup {
	if m == nil {
		return nil
	}
	to := &grid.Upgradegroup{
		Comment:                    flex.ExpandStringPointer(m.Comment),
		DistributionDependentGroup: flex.ExpandStringPointer(m.DistributionDependentGroup),
		DistributionPolicy:         flex.ExpandStringPointer(m.DistributionPolicy),
		Members:                    flex.ExpandFrameworkListNestedBlock(ctx, m.Members, diags, ExpandUpgradegroupMembers),
		Name:                       flex.ExpandStringPointer(m.Name),
		UpgradeDependentGroup:      flex.ExpandStringPointer(m.UpgradeDependentGroup),
		UpgradePolicy:              flex.ExpandStringPointer(m.UpgradePolicy),
	}

	to.DistributionTime = flex.ExpandTimeToUnix(m.DistributionTime, diags)
	to.UpgradeTime = flex.ExpandTimeToUnix(m.UpgradeTime, diags)

	return to
}

func FlattenUpgradegroup(ctx context.Context, from *grid.Upgradegroup, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(UpgradegroupAttrTypes)
	}
	m := UpgradegroupModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, UpgradegroupAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *UpgradegroupModel) Flatten(ctx context.Context, from *grid.Upgradegroup, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = UpgradegroupModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.DistributionDependentGroup = flex.FlattenStringPointer(from.DistributionDependentGroup)
	m.DistributionPolicy = flex.FlattenStringPointer(from.DistributionPolicy)
	m.DistributionTime = flex.FlattenUnixTime(from.DistributionTime, diags)
	m.Members = flex.FlattenFrameworkListNestedBlock(ctx, from.Members, UpgradegroupMembersAttrTypes, diags, FlattenUpgradegroupMembers)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.TimeZone = flex.FlattenStringPointer(from.TimeZone)
	m.UpgradeDependentGroup = flex.FlattenStringPointer(from.UpgradeDependentGroup)
	m.UpgradePolicy = flex.FlattenStringPointer(from.UpgradePolicy)
	m.UpgradeTime = flex.FlattenUnixTime(from.UpgradeTime, diags)
}

func (m *UpgradegroupModel) PutExpand(to *grid.Upgradegroup) *grid.Upgradegroup {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range UpgradegroupResourceSchemaAttributes {
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
