package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type MemberPreProvisioningModel struct {
	HardwareInfo types.List `tfsdk:"hardware_info"`
	Licenses     types.List `tfsdk:"licenses"`
}

var MemberPreProvisioningAttrTypes = map[string]attr.Type{
	"hardware_info": types.ListType{ElemType: types.ObjectType{AttrTypes: MemberpreprovisioningHardwareInfoAttrTypes}},
	"licenses":      types.ListType{ElemType: types.StringType},
}

var MemberPreProvisioningResourceSchemaAttributes = map[string]schema.Attribute{
	"hardware_info": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MemberpreprovisioningHardwareInfoResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "An array of structures that describe the hardware being pre-provisioned.",
	},
	"licenses": schema.ListAttribute{
		ElementType: types.StringType,
		Required:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			customvalidator.StringsInSlice([]string{"cloud_api", "dhcp", "dns", "dtc", "enterprise", "fireeye", "ms_management", "nios", "rpz", "sw_tp", "tp_sub", "vnios"}),
		},
		MarkdownDescription: "An array of license types the pre-provisioned member should have in order to join the Grid, or the licenses that must be allocated to the member when it joins the Grid using the token-based authentication.",
	},
}

func ExpandMemberPreProvisioning(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberPreProvisioning {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberPreProvisioningModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberPreProvisioningModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberPreProvisioning {
	if m == nil {
		return nil
	}
	to := &grid.MemberPreProvisioning{
		HardwareInfo: flex.ExpandFrameworkListNestedBlock(ctx, m.HardwareInfo, diags, ExpandMemberpreprovisioningHardwareInfo),
		Licenses:     flex.ExpandFrameworkListString(ctx, m.Licenses, diags),
	}
	return to
}

func FlattenMemberPreProvisioning(ctx context.Context, from *grid.MemberPreProvisioning, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberPreProvisioningAttrTypes)
	}
	m := MemberPreProvisioningModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberPreProvisioningAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberPreProvisioningModel) Flatten(ctx context.Context, from *grid.MemberPreProvisioning, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberPreProvisioningModel{}
	}
	m.HardwareInfo = flex.FlattenFrameworkListNestedBlock(ctx, from.HardwareInfo, MemberpreprovisioningHardwareInfoAttrTypes, diags, FlattenMemberpreprovisioningHardwareInfo)
	m.Licenses = flex.FlattenFrameworkListString(ctx, from.Licenses, diags)
}

func (m *MemberPreProvisioningModel) PutExpand(to *grid.MemberPreProvisioning) *grid.MemberPreProvisioning {
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

	for field, attr := range MemberPreProvisioningResourceSchemaAttributes {
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
