package security

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/security"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupDnsToplevelCommandsModel struct {
	DdnsAdd          types.Bool `tfsdk:"ddns_add"`
	DdnsDelete       types.Bool `tfsdk:"ddns_delete"`
	Delete           types.Bool `tfsdk:"delete"`
	DnsARecordDelete types.Bool `tfsdk:"dns_a_record_delete"`
	EnableAll        types.Bool `tfsdk:"enable_all"`
	DisableAll       types.Bool `tfsdk:"disable_all"`
}

var AdmingroupDnsToplevelCommandsAttrTypes = map[string]attr.Type{
	"ddns_add":            types.BoolType,
	"ddns_delete":         types.BoolType,
	"delete":              types.BoolType,
	"dns_a_record_delete": types.BoolType,
	"enable_all":          types.BoolType,
	"disable_all":         types.BoolType,
}

var AdmingroupDnsToplevelCommandsResourceSchemaAttributes = map[string]schema.Attribute{
	"ddns_add": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"ddns_delete": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"delete": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"dns_a_record_delete": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"enable_all": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then enable all fields",
	},
	"disable_all": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "If True then disable all fields",
	},
}

func ExpandAdmingroupDnsToplevelCommands(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupDnsToplevelCommands {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupDnsToplevelCommandsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupDnsToplevelCommandsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupDnsToplevelCommands {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupDnsToplevelCommands{
		DdnsAdd:          flex.ExpandBoolPointer(m.DdnsAdd),
		DdnsDelete:       flex.ExpandBoolPointer(m.DdnsDelete),
		Delete:           flex.ExpandBoolPointer(m.Delete),
		DnsARecordDelete: flex.ExpandBoolPointer(m.DnsARecordDelete),
	}
	return to
}

func FlattenAdmingroupDnsToplevelCommands(ctx context.Context, from *security.AdmingroupDnsToplevelCommands, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupDnsToplevelCommandsAttrTypes)
	}
	m := AdmingroupDnsToplevelCommandsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupDnsToplevelCommandsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupDnsToplevelCommandsModel) Flatten(ctx context.Context, from *security.AdmingroupDnsToplevelCommands, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupDnsToplevelCommandsModel{}
	}
	m.DdnsAdd = types.BoolPointerValue(from.DdnsAdd)
	m.DdnsDelete = types.BoolPointerValue(from.DdnsDelete)
	m.Delete = types.BoolPointerValue(from.Delete)
	m.DnsARecordDelete = types.BoolPointerValue(from.DnsARecordDelete)
	m.EnableAll = types.BoolPointerValue(from.EnableAll)
	m.DisableAll = types.BoolPointerValue(from.DisableAll)
}

func (m *AdmingroupDnsToplevelCommandsModel) PutExpand(to *security.AdmingroupDnsToplevelCommands) *security.AdmingroupDnsToplevelCommands {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AdmingroupDnsToplevelCommandsResourceSchemaAttributes {
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
