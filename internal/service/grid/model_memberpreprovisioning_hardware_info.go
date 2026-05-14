package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberpreprovisioningHardwareInfoModel struct {
	Hwtype types.String `tfsdk:"hwtype"`
}

var MemberpreprovisioningHardwareInfoAttrTypes = map[string]attr.Type{
	"hwtype": types.StringType,
}

var MemberpreprovisioningHardwareInfoResourceSchemaAttributes = map[string]schema.Attribute{
	"hwtype": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("CP-V1405", "CP-V2205", "CP-V805", "IB-1415", "IB-1425", "IB-1516", "IB-1526", "IB-2215", "IB-2225", "IB-2326", "IB-4015", "IB-4025", "IB-4126", "IB-815", "IB-825", "IB-926", "IB-FLEX", "IB-V1415", "IB-V1425", "IB-V1516", "IB-V1526", "IB-V2215", "IB-V2225", "IB-V2326", "IB-V4015", "IB-V4025", "IB-V4126", "IB-V815", "IB-V825", "IB-V926"),
		},
		MarkdownDescription: "Hardware type.",
	},
}

func ExpandMemberpreprovisioningHardwareInfo(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberpreprovisioningHardwareInfo {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberpreprovisioningHardwareInfoModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberpreprovisioningHardwareInfoModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberpreprovisioningHardwareInfo {
	if m == nil {
		return nil
	}
	to := &grid.MemberpreprovisioningHardwareInfo{
		Hwtype: flex.ExpandStringPointer(m.Hwtype),
	}
	return to
}

func FlattenMemberpreprovisioningHardwareInfo(ctx context.Context, from *grid.MemberpreprovisioningHardwareInfo, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberpreprovisioningHardwareInfoAttrTypes)
	}
	m := MemberpreprovisioningHardwareInfoModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberpreprovisioningHardwareInfoAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberpreprovisioningHardwareInfoModel) Flatten(ctx context.Context, from *grid.MemberpreprovisioningHardwareInfo, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberpreprovisioningHardwareInfoModel{}
	}
	m.Hwtype = flex.FlattenStringPointer(from.Hwtype)
}

func (m *MemberpreprovisioningHardwareInfoModel) PutExpand(to *grid.MemberpreprovisioningHardwareInfo) *grid.MemberpreprovisioningHardwareInfo {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberpreprovisioningHardwareInfoResourceSchemaAttributes {
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
