package acl

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/acl"

	"github.com/infobloxopen/terraform-provider-nios/internal/flex"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type NamedaclModel struct {
	Ref                types.String `tfsdk:"ref"`
	AccessList         types.List   `tfsdk:"access_list"`
	Comment            types.String `tfsdk:"comment"`
	ExplodedAccessList types.List   `tfsdk:"exploded_access_list"`
	ExtAttrs           types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll        types.Map    `tfsdk:"extattrs_all"`
	Name               types.String `tfsdk:"name"`
}

var NamedaclAttrTypes = map[string]attr.Type{
	"ref":                  types.StringType,
	"access_list":          types.ListType{ElemType: types.ObjectType{AttrTypes: NamedaclAccessListAttrTypes}},
	"comment":              types.StringType,
	"exploded_access_list": types.ListType{ElemType: types.ObjectType{AttrTypes: NamedaclExplodedAccessListAttrTypes}},
	"extattrs":             types.MapType{ElemType: types.StringType},
	"extattrs_all":         types.MapType{ElemType: types.StringType},
	"name":                 types.StringType,
}

var NamedaclResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"access_list": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: NamedaclAccessListResourceSchemaAttributes,
		},
		Optional: true,
		Computed: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Default:             listdefault.StaticValue(types.ListNull(types.ObjectType{AttrTypes: NamedaclAccessListAttrTypes})),
		MarkdownDescription: "The access control list of IPv4/IPv6 addresses, networks, TSIG-based anonymous access controls, and other named ACLs.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthBetween(0, 256),
		},
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "Comment for the named ACL; maximum 256 characters.",
	},
	"exploded_access_list": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: NamedaclExplodedAccessListResourceSchemaAttributes,
		},
		Computed:            true,
		MarkdownDescription: "The exploded access list for the named ACL. This list displays all the access control entries in a named ACL and its nested named ACLs, if applicable.",
	},
	"extattrs": schema.MapAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Extensible attributes associated with the object.",
		ElementType:         types.StringType,
		Default:             mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
	},
	"extattrs_all": schema.MapAttribute{
		Computed:            true,
		MarkdownDescription: "Extensible attributes associated with the object, including default attributes.",
		ElementType:         types.StringType,
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of the named ACL.",
	},
}

func (m *NamedaclModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *acl.Namedacl {
	if m == nil {
		return nil
	}
	to := &acl.Namedacl{
		AccessList: flex.ExpandFrameworkListNestedBlock(ctx, m.AccessList, diags, ExpandNamedaclAccessList),
		Comment:    flex.ExpandStringPointer(m.Comment),
		ExtAttrs:   ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Name:       flex.ExpandStringPointer(m.Name),
	}
	return to
}

func FlattenNamedacl(ctx context.Context, from *acl.Namedacl, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NamedaclAttrTypes)
	}
	m := NamedaclModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, NamedaclAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NamedaclModel) Flatten(ctx context.Context, from *acl.Namedacl, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NamedaclModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AccessList = flattenAccessListWithPlanAddress(ctx, m.AccessList, from.AccessList, diags)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.ExplodedAccessList = flex.FlattenFrameworkListNestedBlock(ctx, from.ExplodedAccessList, NamedaclExplodedAccessListAttrTypes, diags, FlattenNamedaclExplodedAccessList)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Name = flex.FlattenStringPointer(from.Name)
}

// flattenAccessListWithPlanAddress flattens the API access list and reconciles
// each entry's address with the corresponding plan/state value via
// FlattenNamedaclAddress so that user-specified "/32" CIDR suffixes stripped
// by the API are preserved in state.
func flattenAccessListWithPlanAddress(ctx context.Context, planList types.List, fromList []acl.NamedaclAccessList, diags *diag.Diagnostics) types.List {
	var planAccessList []NamedaclAccessListModel
	if !planList.IsNull() && !planList.IsUnknown() {
		diags.Append(planList.ElementsAs(ctx, &planAccessList, false)...)
	}

	apiList := flex.FlattenFrameworkListNestedBlock(ctx, fromList, NamedaclAccessListAttrTypes, diags, FlattenNamedaclAccessList)

	if len(planAccessList) == 0 || apiList.IsNull() || apiList.IsUnknown() {
		return apiList
	}

	var apiAccessList []NamedaclAccessListModel
	diags.Append(apiList.ElementsAs(ctx, &apiAccessList, false)...)
	if diags.HasError() {
		return apiList
	}

	for i := range apiAccessList {
		if i >= len(planAccessList) {
			break
		}
		apiAccessList[i].Address = FlattenNamedaclAddress(planAccessList[i].Address, apiAccessList[i].Address)
	}

	newList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NamedaclAccessListAttrTypes}, apiAccessList)
	diags.Append(d...)
	if d.HasError() {
		return apiList
	}
	return newList
}

// FlattenNamedaclAddress reconciles the API response address with the plan
// address. If the plan address was specified with a "/32" CIDR suffix but the
// API response strips it (returning just the bare IP), the suffix is added
// back so the state matches what the user configured and avoids drift.
func FlattenNamedaclAddress(planAddr, apiAddr types.String) types.String {
	if apiAddr.IsNull() || apiAddr.IsUnknown() {
		return apiAddr
	}
	addr := apiAddr.ValueString()
	if !planAddr.IsNull() && !planAddr.IsUnknown() {
		plan := planAddr.ValueString()
		if strings.HasSuffix(plan, "/32") && !strings.HasSuffix(addr, "/32") {
			addr = addr + "/32"
		}
	}
	return types.StringValue(addr)
}

func (m *NamedaclModel) PutExpand(to *acl.Namedacl) *acl.Namedacl {
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

	for field, attr := range NamedaclResourceSchemaAttributes {
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
