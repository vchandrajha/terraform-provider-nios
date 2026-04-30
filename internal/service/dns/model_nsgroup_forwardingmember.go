package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type NsgroupForwardingmemberModel struct {
	Ref               types.String `tfsdk:"ref"`
	Comment           types.String `tfsdk:"comment"`
	ExtAttrs          types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll       types.Map    `tfsdk:"extattrs_all"`
	ForwardingServers types.List   `tfsdk:"forwarding_servers"`
	Name              types.String `tfsdk:"name"`
}

var NsgroupForwardingmemberAttrTypes = map[string]attr.Type{
	"ref":                types.StringType,
	"comment":            types.StringType,
	"extattrs":           types.MapType{ElemType: types.StringType},
	"extattrs_all":       types.MapType{ElemType: types.StringType},
	"forwarding_servers": types.ListType{ElemType: types.ObjectType{AttrTypes: NsgroupForwardingmemberForwardingServersAttrTypes}},
	"name":               types.StringType,
}

var NsgroupForwardingmemberResourceSchemaAttributes = map[string]schema.Attribute{
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
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Comment for the Forwarding Member Name Server Group; maximum 256 characters.",
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
		MarkdownDescription: "Extensible attributes associated with the object , including default attributes.",
		ElementType:         types.StringType,
		PlanModifiers: []planmodifier.Map{
			importmod.AssociateInternalId(),
			mapplanmodifier.UseStateForUnknown(),
		},
	},
	"forwarding_servers": schema.ListNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: NsgroupForwardingmemberForwardingServersResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The list of forwarding member servers.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of the Forwarding Member Name Server Group.",
	},
}

func (m *NsgroupForwardingmemberModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.NsgroupForwardingmember {
	if m == nil {
		return nil
	}
	to := &dns.NsgroupForwardingmember{
		Comment:           flex.ExpandStringPointer(m.Comment),
		ExtAttrs:          ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		ForwardingServers: flex.ExpandFrameworkListNestedBlock(ctx, m.ForwardingServers, diags, ExpandNsgroupForwardingmemberForwardingServers),
		Name:              flex.ExpandStringPointer(m.Name),
	}
	return to
}

func FlattenNsgroupForwardingmember(ctx context.Context, from *dns.NsgroupForwardingmember, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(NsgroupForwardingmemberAttrTypes)
	}
	m := NsgroupForwardingmemberModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, NsgroupForwardingmemberAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *NsgroupForwardingmemberModel) Flatten(ctx context.Context, from *dns.NsgroupForwardingmember, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = NsgroupForwardingmemberModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.ForwardingServers = flex.FlattenFrameworkListNestedBlock(ctx, from.ForwardingServers, NsgroupForwardingmemberForwardingServersAttrTypes, diags, FlattenNsgroupForwardingmemberForwardingServers)
	m.Name = flex.FlattenStringPointer(from.Name)
}

func (m *NsgroupForwardingmemberModel) PutExpand(to *dns.NsgroupForwardingmember) *dns.NsgroupForwardingmember {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range NsgroupForwardingmemberResourceSchemaAttributes {
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
