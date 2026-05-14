package ipam

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	internaltypes "github.com/infobloxopen/terraform-provider-nios/internal/types"

	"github.com/infobloxopen/infoblox-nios-go-client/ipam"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type SuperhostModel struct {
	Ref                     types.String                     `tfsdk:"ref"`
	Comment                 types.String                     `tfsdk:"comment"`
	DeleteAssociatedObjects types.Bool                       `tfsdk:"delete_associated_objects"`
	DhcpAssociatedObjects   internaltypes.UnorderedListValue `tfsdk:"dhcp_associated_objects"`
	Disabled                types.Bool                       `tfsdk:"disabled"`
	DnsAssociatedObjects    internaltypes.UnorderedListValue `tfsdk:"dns_associated_objects"`
	ExtAttrs                types.Map                        `tfsdk:"extattrs"`
	ExtAttrsAll             types.Map                        `tfsdk:"extattrs_all"`
	Name                    types.String                     `tfsdk:"name"`
}

var SuperhostAttrTypes = map[string]attr.Type{
	"ref":                       types.StringType,
	"comment":                   types.StringType,
	"delete_associated_objects": types.BoolType,
	"dhcp_associated_objects":   internaltypes.UnorderedListOfStringType,
	"disabled":                  types.BoolType,
	"dns_associated_objects":    internaltypes.UnorderedListOfStringType,
	"extattrs":                  types.MapType{ElemType: types.StringType},
	"extattrs_all":              types.MapType{ElemType: types.StringType},
	"name":                      types.StringType,
}

var SuperhostResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The comment for Super Host.",
	},
	"delete_associated_objects": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "True if we have to delete all DNS/DHCP associated objects with Super Host, false by default.",
	},
	"dhcp_associated_objects": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of DHCP objects refs which are associated with Super Host.",
	},
	"disabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Disable all DNS/DHCP associated objects with Super Host if True, False by default.",
	},
	"dns_associated_objects": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of object refs of the DNS resource records which are associated with Super Host.",
	},
	"extattrs": schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Extensible attributes associated with the object.",
	},
	"extattrs_all": schema.MapAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Extensible attributes associated with the object, including default and internal attributes.",
		ElementType:         types.StringType,
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Name of the Super Host.",
	},
}

func (m *SuperhostModel) Expand(ctx context.Context, diags *diag.Diagnostics) *ipam.Superhost {
	if m == nil {
		return nil
	}
	to := &ipam.Superhost{
		Comment:                 flex.ExpandStringPointer(m.Comment),
		DeleteAssociatedObjects: flex.ExpandBoolPointer(m.DeleteAssociatedObjects),
		DhcpAssociatedObjects:   flex.ExpandFrameworkListString(ctx, m.DhcpAssociatedObjects, diags),
		Disabled:                flex.ExpandBoolPointer(m.Disabled),
		DnsAssociatedObjects:    flex.ExpandFrameworkListString(ctx, m.DnsAssociatedObjects, diags),
		ExtAttrs:                ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Name:                    flex.ExpandStringPointer(m.Name),
	}
	return to
}

func FlattenSuperhost(ctx context.Context, from *ipam.Superhost, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SuperhostAttrTypes)
	}
	m := SuperhostModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, SuperhostAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SuperhostModel) Flatten(ctx context.Context, from *ipam.Superhost, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SuperhostModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.DhcpAssociatedObjects = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.DhcpAssociatedObjects, diags)
	m.Disabled = types.BoolPointerValue(from.Disabled)
	m.DnsAssociatedObjects = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.DnsAssociatedObjects, diags)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Name = flex.FlattenStringPointer(from.Name)
}

func (m *SuperhostModel) PutExpand(to *ipam.Superhost) *ipam.Superhost {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range SuperhostResourceSchemaAttributes {
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
