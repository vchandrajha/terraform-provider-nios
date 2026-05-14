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
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type AdmingroupGridSetCommandsModel struct {
	SetDefaultRevertWindow types.Bool `tfsdk:"set_default_revert_window"`
	SetDscp                types.Bool `tfsdk:"set_dscp"`
	SetMembership          types.Bool `tfsdk:"set_membership"`
	SetNogrid              types.Bool `tfsdk:"set_nogrid"`
	SetNomastergrid        types.Bool `tfsdk:"set_nomastergrid"`
	SetPromoteMaster       types.Bool `tfsdk:"set_promote_master"`
	SetRevertGrid          types.Bool `tfsdk:"set_revert_grid"`
	SetToken               types.Bool `tfsdk:"set_token"`
	SetTestPromoteMaster   types.Bool `tfsdk:"set_test_promote_master"`
	EnableAll              types.Bool `tfsdk:"enable_all"`
	DisableAll             types.Bool `tfsdk:"disable_all"`
}

var AdmingroupGridSetCommandsAttrTypes = map[string]attr.Type{
	"set_default_revert_window": types.BoolType,
	"set_dscp":                  types.BoolType,
	"set_membership":            types.BoolType,
	"set_nogrid":                types.BoolType,
	"set_nomastergrid":          types.BoolType,
	"set_promote_master":        types.BoolType,
	"set_revert_grid":           types.BoolType,
	"set_token":                 types.BoolType,
	"set_test_promote_master":   types.BoolType,
	"enable_all":                types.BoolType,
	"disable_all":               types.BoolType,
}

var AdmingroupGridSetCommandsResourceSchemaAttributes = map[string]schema.Attribute{
	"set_default_revert_window": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_dscp": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_membership": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_nogrid": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_nomastergrid": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_promote_master": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_revert_grid": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_token": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"set_test_promote_master": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "If True then CLI user has permission to run the command",
	},
	"enable_all": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "If True then enable all fields",
	},
	"disable_all": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "If True then disable all fields",
	},
}

func ExpandAdmingroupGridSetCommands(ctx context.Context, o types.Object, diags *diag.Diagnostics) *security.AdmingroupGridSetCommands {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m AdmingroupGridSetCommandsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *AdmingroupGridSetCommandsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *security.AdmingroupGridSetCommands {
	if m == nil {
		return nil
	}
	to := &security.AdmingroupGridSetCommands{
		SetDefaultRevertWindow: flex.ExpandBoolPointer(m.SetDefaultRevertWindow),
		SetDscp:                flex.ExpandBoolPointer(m.SetDscp),
		SetMembership:          flex.ExpandBoolPointer(m.SetMembership),
		SetNogrid:              flex.ExpandBoolPointer(m.SetNogrid),
		SetNomastergrid:        flex.ExpandBoolPointer(m.SetNomastergrid),
		SetPromoteMaster:       flex.ExpandBoolPointer(m.SetPromoteMaster),
		SetRevertGrid:          flex.ExpandBoolPointer(m.SetRevertGrid),
		SetToken:               flex.ExpandBoolPointer(m.SetToken),
		SetTestPromoteMaster:   flex.ExpandBoolPointer(m.SetTestPromoteMaster),
	}
	return to
}

func FlattenAdmingroupGridSetCommands(ctx context.Context, from *security.AdmingroupGridSetCommands, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AdmingroupGridSetCommandsAttrTypes)
	}
	m := AdmingroupGridSetCommandsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AdmingroupGridSetCommandsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AdmingroupGridSetCommandsModel) Flatten(ctx context.Context, from *security.AdmingroupGridSetCommands, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AdmingroupGridSetCommandsModel{}
	}
	m.SetDefaultRevertWindow = types.BoolPointerValue(from.SetDefaultRevertWindow)
	m.SetDscp = types.BoolPointerValue(from.SetDscp)
	m.SetMembership = types.BoolPointerValue(from.SetMembership)
	m.SetNogrid = types.BoolPointerValue(from.SetNogrid)
	m.SetNomastergrid = types.BoolPointerValue(from.SetNomastergrid)
	m.SetPromoteMaster = types.BoolPointerValue(from.SetPromoteMaster)
	m.SetRevertGrid = types.BoolPointerValue(from.SetRevertGrid)
	m.SetToken = types.BoolPointerValue(from.SetToken)
	m.SetTestPromoteMaster = types.BoolPointerValue(from.SetTestPromoteMaster)
	m.EnableAll = types.BoolPointerValue(from.EnableAll)
	m.DisableAll = types.BoolPointerValue(from.DisableAll)
}

func (m *AdmingroupGridSetCommandsModel) PutExpand(to *security.AdmingroupGridSetCommands) *security.AdmingroupGridSetCommands {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AdmingroupGridSetCommandsResourceSchemaAttributes {
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
