package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	planmodifiers "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/immutable"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	derivedmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/derived"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type RecordNsModel struct {
	Ref              types.String `tfsdk:"ref"`
	Addresses        types.List   `tfsdk:"addresses"`
	CloudInfo        types.Object `tfsdk:"cloud_info"`
	Creator          types.String `tfsdk:"creator"`
	DnsName          types.String `tfsdk:"dns_name"`
	LastQueried      types.Int64  `tfsdk:"last_queried"`
	MsDelegationName types.String `tfsdk:"ms_delegation_name"`
	Name             types.String `tfsdk:"name"`
	Nameserver       types.String `tfsdk:"nameserver"`
	Policy           types.String `tfsdk:"policy"`
	View             types.String `tfsdk:"view"`
	Zone             types.String `tfsdk:"zone"`
}

var RecordNsAttrTypes = map[string]attr.Type{
	"ref":                types.StringType,
	"addresses":          types.ListType{ElemType: types.ObjectType{AttrTypes: RecordNsAddressesAttrTypes}},
	"cloud_info":         types.ObjectType{AttrTypes: RecordNsCloudInfoAttrTypes},
	"creator":            types.StringType,
	"dns_name":           types.StringType,
	"last_queried":       types.Int64Type,
	"ms_delegation_name": types.StringType,
	"name":               types.StringType,
	"nameserver":         types.StringType,
	"policy":             types.StringType,
	"view":               types.StringType,
	"zone":               types.StringType,
}

var RecordNsResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"addresses": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: RecordNsAddressesResourceSchemaAttributes,
		},
		Required: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The list of zone name servers.",
	},
	"cloud_info": schema.SingleNestedAttribute{
		Attributes:          RecordNsCloudInfoResourceSchemaAttributes,
		Computed:            true,
		MarkdownDescription: "The cloud information associated with the record.",
	},
	"creator": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The record creator.",
	},
	"dns_name": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			derivedmod.PunycodeDerivedFrom("name"),
		},
		MarkdownDescription: "The name of the NS record in punycode format.",
	},
	"last_queried": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The time of the last DNS query in Epoch seconds format.",
	},
	"ms_delegation_name": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The MS delegation point name.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.IsValidDomainName(),
		},
		MarkdownDescription: "The name of the NS record in FQDN format. This value can be in unicode format.",
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
	},
	"nameserver": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The domain name of an authoritative server for the redirected zone.",
	},
	"policy": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The host name policy for the record.",
	},
	"view": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("default"),
		MarkdownDescription: "The name of the DNS view in which the record resides. Example: \"external\".",
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
	},
	"zone": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The name of the zone in which the record resides. Example: \"zone.com\". If a view is not specified when searching by zone, the default view is used.",
	},
}

func (m *RecordNsModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *dns.RecordNs {
	if m == nil {
		return nil
	}
	to := &dns.RecordNs{
		Addresses:        flex.ExpandFrameworkListNestedBlock(ctx, m.Addresses, diags, ExpandRecordNsAddresses),
		MsDelegationName: flex.ExpandStringPointer(m.MsDelegationName),
		Nameserver:       flex.ExpandStringPointer(m.Nameserver),
	}
	if isCreate {
		to.Name = flex.ExpandStringPointer(m.Name)
		to.View = flex.ExpandStringPointer(m.View)
	}
	return to
}

func FlattenRecordNs(ctx context.Context, from *dns.RecordNs, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RecordNsAttrTypes)
	}
	m := RecordNsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RecordNsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RecordNsModel) Flatten(ctx context.Context, from *dns.RecordNs, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RecordNsModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Addresses = flex.FlattenFrameworkListNestedBlock(ctx, from.Addresses, RecordNsAddressesAttrTypes, diags, FlattenRecordNsAddresses)
	m.CloudInfo = FlattenRecordNsCloudInfo(ctx, from.CloudInfo, diags)
	m.Creator = flex.FlattenStringPointer(from.Creator)
	m.DnsName = flex.FlattenStringPointer(from.DnsName)
	m.LastQueried = flex.FlattenInt64Pointer(from.LastQueried)
	m.MsDelegationName = flex.FlattenStringPointer(from.MsDelegationName)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Nameserver = flex.FlattenStringPointer(from.Nameserver)
	m.Policy = flex.FlattenStringPointer(from.Policy)
	m.View = flex.FlattenStringPointer(from.View)
	m.Zone = flex.FlattenStringPointer(from.Zone)
}

func (m *RecordNsModel) PutExpand(to *dns.RecordNs) *dns.RecordNs {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RecordNsResourceSchemaAttributes {
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
