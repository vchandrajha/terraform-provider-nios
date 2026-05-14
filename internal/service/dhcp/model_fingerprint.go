package dhcp

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	planmodifiers "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/immutable"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type FingerprintModel struct {
	Ref                types.String `tfsdk:"ref"`
	Comment            types.String `tfsdk:"comment"`
	DeviceClass        types.String `tfsdk:"device_class"`
	Disable            types.Bool   `tfsdk:"disable"`
	ExtAttrs           types.Map    `tfsdk:"extattrs"`
	ExtAttrsAll        types.Map    `tfsdk:"extattrs_all"`
	Ipv6OptionSequence types.List   `tfsdk:"ipv6_option_sequence"`
	Name               types.String `tfsdk:"name"`
	OptionSequence     types.List   `tfsdk:"option_sequence"`
	Type               types.String `tfsdk:"type"`
	VendorId           types.List   `tfsdk:"vendor_id"`
}

var FingerprintAttrTypes = map[string]attr.Type{
	"ref":                  types.StringType,
	"comment":              types.StringType,
	"device_class":         types.StringType,
	"disable":              types.BoolType,
	"extattrs":             types.MapType{ElemType: types.StringType},
	"extattrs_all":         types.MapType{ElemType: types.StringType},
	"ipv6_option_sequence": types.ListType{ElemType: types.StringType},
	"name":                 types.StringType,
	"option_sequence":      types.ListType{ElemType: types.StringType},
	"type":                 types.StringType,
	"vendor_id":            types.ListType{ElemType: types.StringType},
}

var FingerprintResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The reference to the object.",
	},
	"comment": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
		},
		MarkdownDescription: "Comment for the Fingerprint; maximum 256 characters.",
	},
	"device_class": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthBetween(0, 256),
		},
		MarkdownDescription: "A class of DHCP Fingerprint object; maximum 256 characters.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the DHCP Fingerprint object is disabled or not.",
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
		},
		MarkdownDescription: "Extensible attributes associated with the object, including default and internal attributes.",
		ElementType:         types.StringType,
	},
	"ipv6_option_sequence": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list (comma separated list) of IPv6 option number sequences of the device or operating system.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Name of the DHCP Fingerprint object.",
	},
	"option_sequence": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AtLeastOneOf(
				path.MatchRoot("option_sequence"),
				path.MatchRoot("ipv6_option_sequence"),
				path.MatchRoot("vendor_id"),
			),
		},
		MarkdownDescription: "A list (comma separated list) of IPv4 option number sequences of the device or operating system.",
	},
	"type": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("CUSTOM"),
		Validators: []validator.String{
			stringvalidator.OneOf("CUSTOM"),
		},
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
		MarkdownDescription: "The type of the DHCP Fingerprint object.",
	},
	"vendor_id": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "A list of vendor IDs of the device or operating system.",
	},
}

func (m *FingerprintModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *dhcp.Fingerprint {
	if m == nil {
		return nil
	}
	to := &dhcp.Fingerprint{
		Comment:            flex.ExpandStringPointer(m.Comment),
		DeviceClass:        flex.ExpandStringPointer(m.DeviceClass),
		Disable:            flex.ExpandBoolPointer(m.Disable),
		ExtAttrs:           ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Ipv6OptionSequence: flex.ExpandFrameworkListString(ctx, m.Ipv6OptionSequence, diags),
		Name:               flex.ExpandStringPointer(m.Name),
		OptionSequence:     flex.ExpandFrameworkListString(ctx, m.OptionSequence, diags),
		VendorId:           flex.ExpandFrameworkListString(ctx, m.VendorId, diags),
	}
	if isCreate {
		to.Type = flex.ExpandStringPointer(m.Type)
	}
	return to
}

func FlattenFingerprint(ctx context.Context, from *dhcp.Fingerprint, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(FingerprintAttrTypes)
	}
	m := FingerprintModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, FingerprintAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *FingerprintModel) Flatten(ctx context.Context, from *dhcp.Fingerprint, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = FingerprintModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.DeviceClass = flex.FlattenStringPointer(from.DeviceClass)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Ipv6OptionSequence = flex.FlattenFrameworkListString(ctx, from.Ipv6OptionSequence, diags)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.OptionSequence = flex.FlattenFrameworkListString(ctx, from.OptionSequence, diags)
	m.Type = flex.FlattenStringPointer(from.Type)
	m.VendorId = flex.FlattenFrameworkListString(ctx, from.VendorId, diags)
}

func (m *FingerprintModel) PutExpand(to *dhcp.Fingerprint) *dhcp.Fingerprint {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range FingerprintResourceSchemaAttributes {
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
