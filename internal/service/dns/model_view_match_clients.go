package dns

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type ViewMatchClientsModel struct {
	Ref            types.String `tfsdk:"ref"`
	Address        types.String `tfsdk:"address"`
	Struct         types.String `tfsdk:"struct"`
	Permission     types.String `tfsdk:"permission"`
	TsigKey        types.String `tfsdk:"tsig_key"`
	TsigKeyAlg     types.String `tfsdk:"tsig_key_alg"`
	TsigKeyName    types.String `tfsdk:"tsig_key_name"`
	UseTsigKeyName types.Bool   `tfsdk:"use_tsig_key_name"`
}

var ViewMatchClientsAttrTypes = map[string]attr.Type{
	"ref":               types.StringType,
	"address":           types.StringType,
	"struct":            types.StringType,
	"permission":        types.StringType,
	"tsig_key":          types.StringType,
	"tsig_key_alg":      types.StringType,
	"tsig_key_name":     types.StringType,
	"use_tsig_key_name": types.BoolType,
}

var ViewMatchClientsResourceSchemaAttributes = map[string]schema.Attribute{
	"struct": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("addressac", "tsigac"),
		},
		MarkdownDescription: "The struct type of the object. The value must be one of 'addressac' and 'tsigac'.",
	},
	"ref": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("struct"),
				path.MatchRelative().AtParent().AtName("address"),
				path.MatchRelative().AtParent().AtName("permission"),
				path.MatchRelative().AtParent().AtName("tsig_key"),
				path.MatchRelative().AtParent().AtName("tsig_key_alg"),
				path.MatchRelative().AtParent().AtName("tsig_key_name"),
			),
		},
		MarkdownDescription: "The reference to the Named ACL object.",
	},
	"address": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("tsig_key"),
				path.MatchRelative().AtParent().AtName("tsig_key_alg"),
				path.MatchRelative().AtParent().AtName("use_tsig_key_name"),
			),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The address this rule applies to or \"Any\".",
	},
	"permission": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("tsig_key"),
				path.MatchRelative().AtParent().AtName("tsig_key_alg"),
				path.MatchRelative().AtParent().AtName("use_tsig_key_name"),
			),
		},
		MarkdownDescription: "The permission to use for this address.",
	},
	"tsig_key": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("address"),
				path.MatchRelative().AtParent().AtName("permission"),
			),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "A generated TSIG key. If the external primary server is a NIOS appliance running DNS One 2.x code, this can be set to :2xCOMPAT.",
	},
	"tsig_key_alg": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("address"),
				path.MatchRelative().AtParent().AtName("permission"),
			),
		},
		MarkdownDescription: "The TSIG key algorithm.",
	},
	"tsig_key_name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("address"),
				path.MatchRelative().AtParent().AtName("permission"),
			),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of the TSIG key. If 2.x TSIG compatibility is used, this is set to 'tsig_xfer' on retrieval, and ignored on insert or update.",
	},
	"use_tsig_key_name": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Bool{
			boolvalidator.ConflictsWith(
				path.MatchRelative().AtParent().AtName("address"),
				path.MatchRelative().AtParent().AtName("permission"),
			),
		},
		MarkdownDescription: "Use flag for: tsig_key_name",
	},
}

func ExpandViewMatchClients(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ViewMatchClients {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ViewMatchClientsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ViewMatchClientsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ViewMatchClients {
	if m == nil {
		return nil
	}
	to := &dns.ViewMatchClients{
		Ref:            flex.ExpandStringPointer(m.Ref),
		Address:        flex.ExpandStringPointer(m.Address),
		Struct:         flex.ExpandStringPointer(m.Struct),
		Permission:     flex.ExpandStringPointer(m.Permission),
		TsigKey:        flex.ExpandStringPointer(m.TsigKey),
		TsigKeyAlg:     flex.ExpandStringPointer(m.TsigKeyAlg),
		TsigKeyName:    flex.ExpandStringPointer(m.TsigKeyName),
		UseTsigKeyName: flex.ExpandBoolPointer(m.UseTsigKeyName),
	}
	return to
}

func FlattenViewMatchClients(ctx context.Context, from *dns.ViewMatchClients, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ViewMatchClientsAttrTypes)
	}
	m := ViewMatchClientsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ViewMatchClientsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ViewMatchClientsModel) Flatten(ctx context.Context, from *dns.ViewMatchClients, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ViewMatchClientsModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Address = flex.FlattenStringPointer(from.Address)
	m.Struct = flex.FlattenStringPointer(from.Struct)
	m.Permission = flex.FlattenStringPointer(from.Permission)
	m.TsigKey = flex.FlattenStringPointer(from.TsigKey)
	m.TsigKeyAlg = flex.FlattenStringPointer(from.TsigKeyAlg)
	m.TsigKeyName = flex.FlattenStringPointer(from.TsigKeyName)
	m.UseTsigKeyName = types.BoolPointerValue(from.UseTsigKeyName)
}

func (m *ViewMatchClientsModel) PutExpand(to *dns.ViewMatchClients) *dns.ViewMatchClients {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ViewMatchClientsResourceSchemaAttributes {
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
							fmt.Printf("Field: %s, Computed: %v, fieldValue: %v, Value: %s\n", field, boolComp, fieldValue, txtFieldValue)
							if ok {
								if !boolComp {
									continue
								} else if txtFieldValue == "" {
									utils.DeleteBy(to, tField.Name)
								}
							} else if txtFieldValue == "" {
								fmt.Printf("Field: %s is marked as computed but is not a bool. Value: %s\n", field, txtFieldValue)
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
