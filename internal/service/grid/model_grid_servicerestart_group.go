package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type GridServicerestartGroupModel struct {
	Ref               types.String `tfsdk:"ref"`
	Comment           types.String `tfsdk:"comment"`
	ExtAttrs          types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll       types.Map    `tfsdk:"extattrs_all"`
	IsDefault         types.Bool   `tfsdk:"is_default"`
	LastUpdatedTime   types.Int64  `tfsdk:"last_updated_time"`
	Members           types.List   `tfsdk:"members"`
	Mode              types.String `tfsdk:"mode"`
	Name              types.String `tfsdk:"name"`
	Position          types.Int64  `tfsdk:"position"`
	RecurringSchedule types.Object `tfsdk:"recurring_schedule"`
	Requests          types.List   `tfsdk:"requests"`
	Service           types.String `tfsdk:"service"`
	Status            types.Object `tfsdk:"status"`
}

var GridServicerestartGroupAttrTypes = map[string]attr.Type{
	"ref":                types.StringType,
	"comment":            types.StringType,
	"extattrs":           types.MapType{ElemType: types.StringType},
	"extattrs_all":       types.MapType{ElemType: types.StringType},
	"is_default":         types.BoolType,
	"last_updated_time":  types.Int64Type,
	"members":            types.ListType{ElemType: types.StringType},
	"mode":               types.StringType,
	"name":               types.StringType,
	"position":           types.Int64Type,
	"recurring_schedule": types.ObjectType{AttrTypes: GridServicerestartGroupRecurringScheduleAttrTypes},
	"requests":           types.ListType{ElemType: types.StringType},
	"service":            types.StringType,
	"status":             types.ObjectType{AttrTypes: GridServicerestartStatusAttrTypes},
}

var GridServicerestartGroupResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthBetween(0, 256),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "Comment for the Restart Group; maximum 256 characters.",
	},
	"extattrs": schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Extensible attributes associated with the object. For valid values for extensible attributes, see {extattrs:values}.",
	},
	"extattrs_all": schema.MapAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		MarkdownDescription: "Extensible attributes associated with the object. For valid values for extensible attributes, see {extattrs:values}.",
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
	},
	"is_default": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if this Restart Group is the default group.",
	},
	"last_updated_time": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The timestamp when the status of the latest request has changed.",
	},
	"members": schema.ListAttribute{
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The list of members belonging to the group.",
	},
	"mode": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("SEQUENTIAL", "SIMULTANEOUS"),
		},
		Default:             stringdefault.StaticString("SIMULTANEOUS"),
		MarkdownDescription: "The default restart method for this Restart Group.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of this Restart Group.",
	},
	"position": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The order to restart.",
	},
	"recurring_schedule": schema.SingleNestedAttribute{
		Attributes: GridServicerestartGroupRecurringScheduleResourceSchemaAttributes,
		Optional:   true,
		Computed:   true,
	},
	"requests": schema.ListAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The list of requests associated with a restart group.",
	},
	"service": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("DHCP", "DNS"),
		},
		MarkdownDescription: "The applicable service for this Restart Group.",
	},
	"status": schema.SingleNestedAttribute{
		Attributes:          GridServicerestartStatusResourceSchemaAttributes,
		Computed:            true,
		MarkdownDescription: "The restart status for a restart group.",
	},
}

func (m *GridServicerestartGroupModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.GridServicerestartGroup {
	if m == nil {
		return nil
	}
	to := &grid.GridServicerestartGroup{
		Comment:           flex.ExpandStringPointer(m.Comment),
		ExtAttrs:          ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Members:           flex.ExpandFrameworkListString(ctx, m.Members, diags),
		Mode:              flex.ExpandStringPointer(m.Mode),
		Name:              flex.ExpandStringPointer(m.Name),
		RecurringSchedule: ExpandGridServicerestartGroupRecurringSchedule(ctx, m.RecurringSchedule, diags),
		Service:           flex.ExpandStringPointer(m.Service),
	}
	return to
}

func FlattenGridServicerestartGroup(ctx context.Context, from *grid.GridServicerestartGroup, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(GridServicerestartGroupAttrTypes)
	}
	m := GridServicerestartGroupModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, GridServicerestartGroupAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *GridServicerestartGroupModel) Flatten(ctx context.Context, from *grid.GridServicerestartGroup, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = GridServicerestartGroupModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.IsDefault = types.BoolPointerValue(from.IsDefault)
	m.LastUpdatedTime = flex.FlattenInt64Pointer(from.LastUpdatedTime)
	m.Members = flex.FlattenFrameworkListString(ctx, from.Members, diags)
	m.Mode = flex.FlattenStringPointer(from.Mode)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Position = flex.FlattenInt64Pointer(from.Position)
	m.RecurringSchedule = FlattenGridServicerestartGroupRecurringSchedule(ctx, from.RecurringSchedule, diags)
	m.Requests = flex.FlattenFrameworkListString(ctx, from.Requests, diags)
	m.Service = flex.FlattenStringPointer(from.Service)
	m.Status = FlattenGridServicerestartStatus(ctx, from.Status, diags)
}

func (m *GridServicerestartGroupModel) PutExpand(to *grid.GridServicerestartGroup) *grid.GridServicerestartGroup {
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

	for field, attr := range GridServicerestartGroupResourceSchemaAttributes {
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
