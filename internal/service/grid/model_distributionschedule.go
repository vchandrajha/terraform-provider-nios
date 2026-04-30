package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type DistributionscheduleModel struct {
	Ref           types.String `tfsdk:"ref"`
	Active        types.Bool   `tfsdk:"active"`
	StartTime     types.String `tfsdk:"start_time"`
	TimeZone      types.String `tfsdk:"time_zone"`
	UpgradeGroups types.List   `tfsdk:"upgrade_groups"`
}

var DistributionscheduleAttrTypes = map[string]attr.Type{
	"ref":            types.StringType,
	"active":         types.BoolType,
	"start_time":     types.StringType,
	"time_zone":      types.StringType,
	"upgrade_groups": types.ListType{ElemType: types.ObjectType{AttrTypes: DistributionscheduleUpgradeGroupsAttrTypes}},
}

var DistributionscheduleResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"active": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines whether the distribution schedule is active.",
	},
	"start_time": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The start time of the distribution.",
		Validators: []validator.String{
			customvalidator.ValidateTimeFormat(),
		},
	},
	"time_zone": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Time zone of the distribution start time.",
	},
	"upgrade_groups": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: DistributionscheduleUpgradeGroupsResourceSchemaAttributes,
		},
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The upgrade groups scheduling settings.",
	},
}

func (m *DistributionscheduleModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.Distributionschedule {
	var groups []grid.DistributionscheduleUpgradeGroups

	if m == nil {
		return nil
	}

	allGroups := flex.ExpandFrameworkListNestedBlock(ctx, m.UpgradeGroups, diags, ExpandDistributionscheduleUpgradeGroups)

	for _, group := range allGroups {
		// Convert empty optional fields to nil
		if group.UpgradeDependentGroup != nil && *group.UpgradeDependentGroup == "" {
			group.UpgradeDependentGroup = nil
		}
		if group.DistributionDependentGroup != nil && *group.DistributionDependentGroup == "" {
			group.DistributionDependentGroup = nil
		}

		// UpgradeTime cannot be nil, set to 0 if not provided
		if group.UpgradeTime == nil {
			val := int64(0)
			group.UpgradeTime = &val
		}

		groups = append(groups, group)
	}

	to := &grid.Distributionschedule{
		Active:        flex.ExpandBoolPointer(m.Active),
		UpgradeGroups: groups,
	}

	to.StartTime = flex.ExpandTimeToUnix(m.StartTime, diags)

	return to
}

func FlattenDistributionschedule(ctx context.Context, from *grid.Distributionschedule, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DistributionscheduleAttrTypes)
	}
	m := DistributionscheduleModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, DistributionscheduleAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DistributionscheduleModel) Flatten(ctx context.Context, from *grid.Distributionschedule, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DistributionscheduleModel{}
	}

	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Active = types.BoolPointerValue(from.Active)
	m.StartTime = flex.FlattenUnixTime(from.StartTime, diags)
	m.TimeZone = flex.FlattenStringPointer(from.TimeZone)
	m.UpgradeGroups = flex.FlattenFrameworkListNestedBlock(ctx, from.UpgradeGroups, DistributionscheduleUpgradeGroupsAttrTypes, diags, FlattenDistributionscheduleUpgradeGroups)
}

func (m *DistributionscheduleModel) PutExpand(to *grid.Distributionschedule) *grid.Distributionschedule {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range DistributionscheduleResourceSchemaAttributes {
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
