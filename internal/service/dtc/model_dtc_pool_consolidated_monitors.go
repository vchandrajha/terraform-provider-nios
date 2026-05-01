package dtc

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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dtc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type DtcPoolConsolidatedMonitorsModel struct {
	Members                 types.List   `tfsdk:"members"`
	Monitor                 types.String `tfsdk:"monitor"`
	Availability            types.String `tfsdk:"availability"`
	FullHealthCommunication types.Bool   `tfsdk:"full_health_communication"`
}

var DtcPoolConsolidatedMonitorsAttrTypes = map[string]attr.Type{
	"members":                   types.ListType{ElemType: types.StringType},
	"monitor":                   types.StringType,
	"availability":              types.StringType,
	"full_health_communication": types.BoolType,
}

var DtcPoolConsolidatedMonitorsResourceSchemaAttributes = map[string]schema.Attribute{
	"members": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Members whose monitor statuses are shared across other members in a pool.",
	},
	"monitor": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Monitor whose statuses are shared across other members in a pool.",
	},
	"availability": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("ANY", "ALL"),
		},
		MarkdownDescription: "Servers assigned to a pool with monitor defined are healthy if ANY or ALL members report healthy status.",
	},
	"full_health_communication": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag for switching health performing and sharing behavior to perform health checks on each DTC grid member that serves related LBDN(s) and send them across all DTC grid members from both selected and non-selected lists.",
	},
}

func ExpandDtcPoolConsolidatedMonitors(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dtc.DtcPoolConsolidatedMonitors {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m DtcPoolConsolidatedMonitorsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *DtcPoolConsolidatedMonitorsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dtc.DtcPoolConsolidatedMonitors {
	if m == nil {
		return nil
	}
	to := &dtc.DtcPoolConsolidatedMonitors{
		Members:                 flex.ExpandFrameworkListString(ctx, m.Members, diags),
		Monitor:                 flex.ExpandStringPointer(m.Monitor),
		Availability:            flex.ExpandStringPointer(m.Availability),
		FullHealthCommunication: flex.ExpandBoolPointer(m.FullHealthCommunication),
	}
	return to
}

func FlattenDtcPoolConsolidatedMonitors(ctx context.Context, from *dtc.DtcPoolConsolidatedMonitors, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DtcPoolConsolidatedMonitorsAttrTypes)
	}
	m := DtcPoolConsolidatedMonitorsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, DtcPoolConsolidatedMonitorsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DtcPoolConsolidatedMonitorsModel) Flatten(ctx context.Context, from *dtc.DtcPoolConsolidatedMonitors, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DtcPoolConsolidatedMonitorsModel{}
	}
	m.Members = flex.FlattenFrameworkListString(ctx, from.Members, diags)
	m.Monitor = flex.FlattenStringPointer(from.Monitor)
	m.Availability = flex.FlattenStringPointer(from.Availability)
	m.FullHealthCommunication = types.BoolPointerValue(from.FullHealthCommunication)
}

func (m *DtcPoolConsolidatedMonitorsModel) PutExpand(to *dtc.DtcPoolConsolidatedMonitors) *dtc.DtcPoolConsolidatedMonitors {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range DtcPoolConsolidatedMonitorsResourceSchemaAttributes {
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
