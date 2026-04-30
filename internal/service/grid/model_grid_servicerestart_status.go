package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type GridServicerestartStatusModel struct {
	Ref            types.String `tfsdk:"ref"`
	Failures       types.Int64  `tfsdk:"failures"`
	Finished       types.Int64  `tfsdk:"finished"`
	Grouped        types.String `tfsdk:"grouped"`
	NeededRestart  types.Int64  `tfsdk:"needed_restart"`
	NoRestart      types.Int64  `tfsdk:"no_restart"`
	Parent         types.String `tfsdk:"parent"`
	Pending        types.Int64  `tfsdk:"pending"`
	PendingRestart types.Int64  `tfsdk:"pending_restart"`
	Processing     types.Int64  `tfsdk:"processing"`
	Restarting     types.Int64  `tfsdk:"restarting"`
	Success        types.Int64  `tfsdk:"success"`
	Timeouts       types.Int64  `tfsdk:"timeouts"`
}

var GridServicerestartStatusAttrTypes = map[string]attr.Type{
	"ref":             types.StringType,
	"failures":        types.Int64Type,
	"finished":        types.Int64Type,
	"grouped":         types.StringType,
	"needed_restart":  types.Int64Type,
	"no_restart":      types.Int64Type,
	"parent":          types.StringType,
	"pending":         types.Int64Type,
	"pending_restart": types.Int64Type,
	"processing":      types.Int64Type,
	"restarting":      types.Int64Type,
	"success":         types.Int64Type,
	"timeouts":        types.Int64Type,
}

var GridServicerestartStatusResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"failures": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of failed requests.",
	},
	"finished": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of finished requests.",
	},
	"grouped": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The type of grouping.",
	},
	"needed_restart": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of created yet unprocessed requests for restart.",
	},
	"no_restart": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of requests that did not require a restart.",
	},
	"parent": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "A reference to the grid or grid:servicerestart:group object associated with the request.",
	},
	"pending": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of requests that are pending a restart.",
	},
	"pending_restart": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of forced or needed requests pending for restart.",
	},
	"processing": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of not forced and not needed requests pending for restart.",
	},
	"restarting": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of service restarts for corresponding members.",
	},
	"success": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of requests associated with successful restarts.",
	},
	"timeouts": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The number of timeout requests.",
	},
}

func ExpandGridServicerestartStatus(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.GridServicerestartStatus {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m GridServicerestartStatusModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *GridServicerestartStatusModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.GridServicerestartStatus {
	if m == nil {
		return nil
	}
	to := &grid.GridServicerestartStatus{}
	return to
}

func FlattenGridServicerestartStatus(ctx context.Context, from *grid.GridServicerestartStatus, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(GridServicerestartStatusAttrTypes)
	}
	m := GridServicerestartStatusModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, GridServicerestartStatusAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *GridServicerestartStatusModel) Flatten(ctx context.Context, from *grid.GridServicerestartStatus, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = GridServicerestartStatusModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Failures = flex.FlattenInt64Pointer(from.Failures)
	m.Finished = flex.FlattenInt64Pointer(from.Finished)
	m.Grouped = flex.FlattenStringPointer(from.Grouped)
	m.NeededRestart = flex.FlattenInt64Pointer(from.NeededRestart)
	m.NoRestart = flex.FlattenInt64Pointer(from.NoRestart)
	m.Parent = flex.FlattenStringPointer(from.Parent)
	m.Pending = flex.FlattenInt64Pointer(from.Pending)
	m.PendingRestart = flex.FlattenInt64Pointer(from.PendingRestart)
	m.Processing = flex.FlattenInt64Pointer(from.Processing)
	m.Restarting = flex.FlattenInt64Pointer(from.Restarting)
	m.Success = flex.FlattenInt64Pointer(from.Success)
	m.Timeouts = flex.FlattenInt64Pointer(from.Timeouts)
}

func (m *GridServicerestartStatusModel) PutExpand(to *grid.GridServicerestartStatus) *grid.GridServicerestartStatus {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range GridServicerestartStatusResourceSchemaAttributes {
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
