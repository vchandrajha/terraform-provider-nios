package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	internaltypes "github.com/infobloxopen/terraform-provider-nios/internal/types"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type GridServicerestartGroupRecurringScheduleModel struct {
	Services internaltypes.UnorderedListValue `tfsdk:"services"`
	Mode     types.String                     `tfsdk:"mode"`
	Schedule types.Object                     `tfsdk:"schedule"`
	Force    types.Bool                       `tfsdk:"force"`
}

var GridServicerestartGroupRecurringScheduleAttrTypes = map[string]attr.Type{
	"services": internaltypes.UnorderedListOfStringType,
	"mode":     types.StringType,
	"schedule": types.ObjectType{AttrTypes: GridservicerestartgrouprecurringscheduleScheduleAttrTypes},
	"force":    types.BoolType,
}

var GridServicerestartGroupRecurringScheduleResourceSchemaAttributes = map[string]schema.Attribute{
	"services": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Validators: []validator.List{
			customvalidator.StringsInSlice([]string{"ALL", "DHCP", "DHCPV4", "DHCPV6", "DNS"}),
			listvalidator.SizeAtLeast(1),
		},
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The list of applicable services for the restart.",
	},
	"mode": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("GROUPED", "SEQUENTIAL", "SIMULTANEOUS"),
		},
		MarkdownDescription: "The restart method for a Grid restart.",
	},
	"schedule": schema.SingleNestedAttribute{
		Attributes: GridservicerestartgrouprecurringscheduleScheduleResourceSchemaAttributes,
		Optional:   true,
		Computed:   true,
	},
	"force": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the Restart Group should have a force restart.",
	},
}

func ExpandGridServicerestartGroupRecurringSchedule(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.GridServicerestartGroupRecurringSchedule {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m GridServicerestartGroupRecurringScheduleModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *GridServicerestartGroupRecurringScheduleModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.GridServicerestartGroupRecurringSchedule {
	if m == nil {
		return nil
	}
	to := &grid.GridServicerestartGroupRecurringSchedule{
		Services: flex.ExpandFrameworkListString(ctx, m.Services, diags),
		Mode:     flex.ExpandStringPointer(m.Mode),
		Schedule: ExpandGridservicerestartgrouprecurringscheduleSchedule(ctx, m.Schedule, diags),
		Force:    flex.ExpandBoolPointer(m.Force),
	}
	return to
}

func FlattenGridServicerestartGroupRecurringSchedule(ctx context.Context, from *grid.GridServicerestartGroupRecurringSchedule, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(GridServicerestartGroupRecurringScheduleAttrTypes)
	}
	m := GridServicerestartGroupRecurringScheduleModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, GridServicerestartGroupRecurringScheduleAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *GridServicerestartGroupRecurringScheduleModel) Flatten(ctx context.Context, from *grid.GridServicerestartGroupRecurringSchedule, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = GridServicerestartGroupRecurringScheduleModel{}
	}
	m.Services = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.Services, diags)
	m.Mode = flex.FlattenStringPointer(from.Mode)
	m.Schedule = FlattenGridservicerestartgrouprecurringscheduleSchedule(ctx, from.Schedule, diags)
	m.Force = types.BoolPointerValue(from.Force)
}

func (m *GridServicerestartGroupRecurringScheduleModel) PutExpand(to *grid.GridServicerestartGroupRecurringSchedule) *grid.GridServicerestartGroupRecurringSchedule {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()

	// Helper to recursively delete empty fields in structs
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

	for field, attr := range GridServicerestartGroupRecurringScheduleResourceSchemaAttributes {
		attrVal := reflect.ValueOf(attr)
		attrType := attrVal.Type()
		if toType.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < toType.NumField(); i++ {
			tField := toType.Field(i)
			fieldValue := toVal.Field(i).Interface()
			cleanTag := strings.Split(tField.Tag.Get("json"), ",")[0]
			cleanTag = strings.Trim(cleanTag, "_")
			txtFieldValue := utils.ToString(field, fieldValue)
			if field != cleanTag {
				continue
			}

			// Skip if attribute is Required
			if _, ok := attrType.FieldByName("Required"); ok {
				requiredVal := attrVal.FieldByName("Required")
				if requiredVal.IsValid() && requiredVal.CanInterface() {
					boolReq, ok := requiredVal.Interface().(bool)
					if ok && boolReq {
						continue
					}
				}
			}

			// Handle Default
			if _, ok := attrType.FieldByName("Default"); ok {
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

			// Handle Computed
			if _, ok := attrType.FieldByName("Computed"); ok {
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

			// Recursively clean up nested structs and slices
			fvType := reflect.TypeOf(fieldValue)
			if fvType != nil {
				switch fvType.Kind() {
				case reflect.Struct:
					deleteEmptyFields(reflect.ValueOf(fieldValue))
				case reflect.Slice, reflect.Array:
					sliceVal := reflect.ValueOf(fieldValue)
					for j := 0; j < sliceVal.Len(); j++ {
						elem := sliceVal.Index(j)
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
	return to
}
