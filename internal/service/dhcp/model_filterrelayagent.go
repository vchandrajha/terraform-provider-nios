package dhcp

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/dhcp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type FilterrelayagentModel struct {
	Ref                      types.String `tfsdk:"ref"`
	CircuitIdName            types.String `tfsdk:"circuit_id_name"`
	CircuitIdSubstringLength types.Int64  `tfsdk:"circuit_id_substring_length"`
	CircuitIdSubstringOffset types.Int64  `tfsdk:"circuit_id_substring_offset"`
	Comment                  types.String `tfsdk:"comment"`
	ExtAttrs                 types.Map    `tfsdk:"extattrs"`
	IsCircuitId              types.String `tfsdk:"is_circuit_id"`
	IsCircuitIdSubstring     types.Bool   `tfsdk:"is_circuit_id_substring"`
	IsRemoteId               types.String `tfsdk:"is_remote_id"`
	IsRemoteIdSubstring      types.Bool   `tfsdk:"is_remote_id_substring"`
	Name                     types.String `tfsdk:"name"`
	RemoteIdName             types.String `tfsdk:"remote_id_name"`
	RemoteIdSubstringLength  types.Int64  `tfsdk:"remote_id_substring_length"`
	RemoteIdSubstringOffset  types.Int64  `tfsdk:"remote_id_substring_offset"`
	ExtAttrsAll              types.Map    `tfsdk:"extattrs_all"`
}

var FilterrelayagentAttrTypes = map[string]attr.Type{
	"ref":                         types.StringType,
	"circuit_id_name":             types.StringType,
	"circuit_id_substring_length": types.Int64Type,
	"circuit_id_substring_offset": types.Int64Type,
	"comment":                     types.StringType,
	"extattrs":                    types.MapType{ElemType: types.StringType},
	"is_circuit_id":               types.StringType,
	"is_circuit_id_substring":     types.BoolType,
	"is_remote_id":                types.StringType,
	"is_remote_id_substring":      types.BoolType,
	"name":                        types.StringType,
	"remote_id_name":              types.StringType,
	"remote_id_substring_length":  types.Int64Type,
	"remote_id_substring_offset":  types.Int64Type,
	"extattrs_all":                types.MapType{ElemType: types.StringType},
}

var FilterrelayagentResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"circuit_id_name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The circuit_id_name of a DHCP relay agent filter object. This filter identifies the circuit between the remote host and the relay agent. For example, the identifier can be the ingress interface number of the circuit access unit, perhaps concatenated with the unit ID number and slot number. Also, the circuit ID can be an ATM virtual circuit ID or cable data virtual circuit ID.",
	},
	"circuit_id_substring_length": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The circuit ID substring length.",
	},
	"circuit_id_substring_offset": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The circuit ID substring offset.",
	},
	"comment": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
		},
		MarkdownDescription: "A descriptive comment of a DHCP relay agent filter object.",
	},
	"extattrs": schema.MapAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Default:     mapdefault.StaticValue(types.MapNull(types.StringType)),
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "Extensible attributes associated with the object. For valid values for extensible attributes, see {extattrs:values}.",
	},
	"is_circuit_id": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("MATCHES_VALUE", "ANY", "NOT_SET"),
		},
		Default:             stringdefault.StaticString("ANY"),
		MarkdownDescription: "The circuit ID matching rule of a DHCP relay agent filter object. The circuit_id value takes effect only if the value is \"MATCHES_VALUE\".",
	},
	"is_circuit_id_substring": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "Determines if the substring of circuit ID, instead of the full circuit ID, is matched.",
	},
	"is_remote_id": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("MATCHES_VALUE", "ANY", "NOT_SET"),
		},
		Default:             stringdefault.StaticString("ANY"),
		MarkdownDescription: "The remote ID matching rule of a DHCP relay agent filter object. The remote_id value takes effect only if the value is Matches_Value.",
	},
	"is_remote_id_substring": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "Determines if the substring of remote ID, instead of the full remote ID, is matched.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The name of a DHCP relay agent filter object.",
	},
	"remote_id_name": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The remote ID name attribute of a relay agent filter object. This filter identifies the remote host. The remote ID name can represent many different things such as the caller ID telephone number for a dial-up connection, a user name for logging in to the ISP, a modem ID, etc. When the remote ID name is defined on the relay agent, the DHCP server will have a trusted relationship to identify the remote host. The remote ID name is considered as a trusted identifier.",
	},
	"remote_id_substring_length": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The remote ID substring length.",
	},
	"remote_id_substring_offset": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Optional:            true,
		MarkdownDescription: "The remote ID substring offset.",
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
}

func (m *FilterrelayagentModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dhcp.Filterrelayagent {
	if m == nil {
		return nil
	}
	to := &dhcp.Filterrelayagent{
		CircuitIdName:            flex.ExpandStringPointer(m.CircuitIdName),
		CircuitIdSubstringLength: flex.ExpandInt64Pointer(m.CircuitIdSubstringLength),
		CircuitIdSubstringOffset: flex.ExpandInt64Pointer(m.CircuitIdSubstringOffset),
		Comment:                  flex.ExpandStringPointer(m.Comment),
		ExtAttrs:                 ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		IsCircuitId:              flex.ExpandStringPointer(m.IsCircuitId),
		IsCircuitIdSubstring:     flex.ExpandBoolPointer(m.IsCircuitIdSubstring),
		IsRemoteId:               flex.ExpandStringPointer(m.IsRemoteId),
		IsRemoteIdSubstring:      flex.ExpandBoolPointer(m.IsRemoteIdSubstring),
		Name:                     flex.ExpandStringPointer(m.Name),
		RemoteIdName:             flex.ExpandStringPointer(m.RemoteIdName),
		RemoteIdSubstringLength:  flex.ExpandInt64Pointer(m.RemoteIdSubstringLength),
		RemoteIdSubstringOffset:  flex.ExpandInt64Pointer(m.RemoteIdSubstringOffset),
	}
	return to
}

func FlattenFilterrelayagent(ctx context.Context, from *dhcp.Filterrelayagent, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(FilterrelayagentAttrTypes)
	}
	m := FilterrelayagentModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, FilterrelayagentAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *FilterrelayagentModel) Flatten(ctx context.Context, from *dhcp.Filterrelayagent, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = FilterrelayagentModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.CircuitIdName = flex.FlattenStringPointer(from.CircuitIdName)
	m.CircuitIdSubstringLength = flex.FlattenInt64Pointer(from.CircuitIdSubstringLength)
	m.CircuitIdSubstringOffset = flex.FlattenInt64Pointer(from.CircuitIdSubstringOffset)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.IsCircuitId = flex.FlattenStringPointer(from.IsCircuitId)
	m.IsCircuitIdSubstring = types.BoolPointerValue(from.IsCircuitIdSubstring)
	m.IsRemoteId = flex.FlattenStringPointer(from.IsRemoteId)
	m.IsRemoteIdSubstring = types.BoolPointerValue(from.IsRemoteIdSubstring)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.RemoteIdName = flex.FlattenStringPointer(from.RemoteIdName)
	m.RemoteIdSubstringLength = flex.FlattenInt64Pointer(from.RemoteIdSubstringLength)
	m.RemoteIdSubstringOffset = flex.FlattenInt64Pointer(from.RemoteIdSubstringOffset)
}

func (m *FilterrelayagentModel) PutExpand(to *dhcp.Filterrelayagent) *dhcp.Filterrelayagent {
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

	for field, attr := range FilterrelayagentResourceSchemaAttributes {
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
