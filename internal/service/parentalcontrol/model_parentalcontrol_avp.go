package parentalcontrol

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/parentalcontrol"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	internaltypes "github.com/infobloxopen/terraform-provider-nios/internal/types"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type ParentalcontrolAvpModel struct {
	Ref          types.String                     `tfsdk:"ref"`
	Comment      types.String                     `tfsdk:"comment"`
	DomainTypes  internaltypes.UnorderedListValue `tfsdk:"domain_types"`
	IsRestricted types.Bool                       `tfsdk:"is_restricted"`
	Name         types.String                     `tfsdk:"name"`
	Type         types.Int64                      `tfsdk:"type"`
	UserDefined  types.Bool                       `tfsdk:"user_defined"`
	ValueType    types.String                     `tfsdk:"value_type"`
	VendorId     types.Int64                      `tfsdk:"vendor_id"`
	VendorType   types.Int64                      `tfsdk:"vendor_type"`
}

var ParentalcontrolAvpAttrTypes = map[string]attr.Type{
	"ref":           types.StringType,
	"comment":       types.StringType,
	"domain_types":  internaltypes.UnorderedListOfStringType,
	"is_restricted": types.BoolType,
	"name":          types.StringType,
	"type":          types.Int64Type,
	"user_defined":  types.BoolType,
	"value_type":    types.StringType,
	"vendor_id":     types.Int64Type,
	"vendor_type":   types.Int64Type,
}

var ParentalcontrolAvpResourceSchemaAttributes = map[string]schema.Attribute{
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
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
		},
		MarkdownDescription: "The comment for the AVP.",
	},
	"domain_types": schema.ListAttribute{
		CustomType:  internaltypes.UnorderedListOfStringType,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AlsoRequires(path.MatchRoot("is_restricted")),
			listvalidator.ValueStringsAre(
				stringvalidator.OneOf(
					"ANCILLARY",
					"IP_SPACE_DIS",
					"NAS_CONTEXT",
					"SUBS_ID")),
		},
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The list of domains applicable to AVP.",
	},
	"is_restricted": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if AVP is restricted to domains.",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The name of AVP.",
	},
	"type": schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 255),
		},
		MarkdownDescription: "The type of AVP as per RFC 2865/2866.",
	},
	"user_defined": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines if AVP was defined by user.",
	},
	"value_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("BYTE", "DATE", "INTEGER", "INTEGER64", "IPADDR", "IPV6ADDR", "IPV6IFID", "IPV6PREFIX", "OCTETS", "SHORT", "STRING"),
		},
		MarkdownDescription: "The type of value.",
	},
	"vendor_id": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The vendor ID as per RFC 2865/2866.",
	},
	"vendor_type": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.Between(1, 255),
		},
		MarkdownDescription: "The vendor type as per RFC 2865/2866.",
	},
}

func ExpandParentalcontrolAvp(ctx context.Context, o types.Object, diags *diag.Diagnostics) *parentalcontrol.ParentalcontrolAvp {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ParentalcontrolAvpModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ParentalcontrolAvpModel) Expand(ctx context.Context, diags *diag.Diagnostics) *parentalcontrol.ParentalcontrolAvp {
	if m == nil {
		return nil
	}
	to := &parentalcontrol.ParentalcontrolAvp{
		Comment:      flex.ExpandStringPointer(m.Comment),
		DomainTypes:  flex.ExpandFrameworkListString(ctx, m.DomainTypes, diags),
		IsRestricted: flex.ExpandBoolPointer(m.IsRestricted),
		Name:         flex.ExpandStringPointer(m.Name),
		Type:         flex.ExpandInt64Pointer(m.Type),
		ValueType:    flex.ExpandStringPointer(m.ValueType),
		VendorId:     flex.ExpandInt64Pointer(m.VendorId),
		VendorType:   flex.ExpandInt64Pointer(m.VendorType),
	}
	return to
}

func FlattenParentalcontrolAvp(ctx context.Context, from *parentalcontrol.ParentalcontrolAvp, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ParentalcontrolAvpAttrTypes)
	}
	m := ParentalcontrolAvpModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ParentalcontrolAvpAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ParentalcontrolAvpModel) Flatten(ctx context.Context, from *parentalcontrol.ParentalcontrolAvp, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ParentalcontrolAvpModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.DomainTypes = flex.FlattenFrameworkUnorderedList(ctx, types.StringType, from.DomainTypes, diags)
	m.IsRestricted = types.BoolPointerValue(from.IsRestricted)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Type = flex.FlattenInt64Pointer(from.Type)
	m.UserDefined = types.BoolPointerValue(from.UserDefined)
	m.ValueType = flex.FlattenStringPointer(from.ValueType)
	m.VendorId = flex.FlattenInt64Pointer(from.VendorId)
	m.VendorType = flex.FlattenInt64Pointer(from.VendorType)
}

func (m *ParentalcontrolAvpModel) PutExpand(to *parentalcontrol.ParentalcontrolAvp) *parentalcontrol.ParentalcontrolAvp {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ParentalcontrolAvpResourceSchemaAttributes {
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
