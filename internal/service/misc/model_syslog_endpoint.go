package misc

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/misc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	importmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/import"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type SyslogEndpointModel struct {
	Ref                types.String `tfsdk:"ref"`
	ExtAttrs           types.Map    `tfsdk:"extattrs"`
	LogLevel           types.String `tfsdk:"log_level"`
	Name               types.String `tfsdk:"name"`
	OutboundMemberType types.String `tfsdk:"outbound_member_type"`
	OutboundMembers    types.List   `tfsdk:"outbound_members"`
	SyslogServers      types.List   `tfsdk:"syslog_servers"`
	TemplateInstance   types.Object `tfsdk:"template_instance"`
	Timeout            types.Int64  `tfsdk:"timeout"`
	VendorIdentifier   types.String `tfsdk:"vendor_identifier"`
	WapiUserName       types.String `tfsdk:"wapi_user_name"`
	WapiUserPassword   types.String `tfsdk:"wapi_user_password"`
	ExtAttrsAll        types.Map    `tfsdk:"extattrs_all"`
}

var SyslogEndpointAttrTypes = map[string]attr.Type{
	"ref":                  types.StringType,
	"extattrs":             types.MapType{ElemType: types.StringType},
	"log_level":            types.StringType,
	"name":                 types.StringType,
	"outbound_member_type": types.StringType,
	"outbound_members":     types.ListType{ElemType: types.StringType},
	"syslog_servers":       types.ListType{ElemType: types.ObjectType{AttrTypes: SyslogEndpointSyslogServersAttrTypes}},
	"template_instance":    types.ObjectType{AttrTypes: SyslogEndpointTemplateInstanceAttrTypes},
	"timeout":              types.Int64Type,
	"vendor_identifier":    types.StringType,
	"wapi_user_name":       types.StringType,
	"wapi_user_password":   types.StringType,
	"extattrs_all":         types.MapType{ElemType: types.StringType},
}

var SyslogEndpointResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
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
	"log_level": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("DEBUG", "ERROR", "INFO", "WARNING"),
		},
		Default:             stringdefault.StaticString("WARNING"),
		MarkdownDescription: "The log level for a notification REST endpoint.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateSyslogEndpointName(),
		},
		MarkdownDescription: "The name of a Syslog endpoint.",
	},
	"outbound_member_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("GM", "MEMBER"),
		},
		MarkdownDescription: "The outbound member that will generate events.",
	},
	"outbound_members": schema.ListAttribute{
		ElementType:         types.StringType,
		Computed:            true,
		Optional:            true,
		Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
		MarkdownDescription: "The list of members for outbound events.",
	},
	"syslog_servers": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: SyslogEndpointSyslogServersResourceSchemaAttributes,
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		Optional: true,
		Computed: true,
		Default: listdefault.StaticValue(
			types.ListValueMust(
				types.ObjectType{AttrTypes: SyslogEndpointSyslogServersAttrTypes},
				[]attr.Value{},
			),
		),
		MarkdownDescription: "List of syslog servers",
	},
	"template_instance": schema.SingleNestedAttribute{
		Attributes:          SyslogEndpointTemplateInstanceResourceSchemaAttributes,
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "The Syslog template instance.",
	},
	"timeout": schema.Int64Attribute{
		Computed:            true,
		Optional:            true,
		Default:             int64default.StaticInt64(30),
		MarkdownDescription: "The timeout of session management (in seconds).",
	},
	"vendor_identifier": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The vendor identifier.",
	},
	"wapi_user_name": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString(""),
		MarkdownDescription: "The user name for WAPI integration.",
	},
	"wapi_user_password": schema.StringAttribute{
		Optional:  true,
		Computed:  true,
		Sensitive: true,
		Default:   stringdefault.StaticString(""),
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The user password for WAPI integration.",
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

func (m *SyslogEndpointModel) Expand(ctx context.Context, diags *diag.Diagnostics) *misc.SyslogEndpoint {
	if m == nil {
		return nil
	}
	to := &misc.SyslogEndpoint{
		ExtAttrs:           ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		LogLevel:           flex.ExpandStringPointer(m.LogLevel),
		Name:               flex.ExpandStringPointer(m.Name),
		OutboundMemberType: flex.ExpandStringPointer(m.OutboundMemberType),
		OutboundMembers:    flex.ExpandFrameworkListString(ctx, m.OutboundMembers, diags),
		SyslogServers:      flex.ExpandFrameworkListNestedBlock(ctx, m.SyslogServers, diags, ExpandSyslogEndpointSyslogServers),
		TemplateInstance:   ExpandSyslogEndpointTemplateInstance(ctx, m.TemplateInstance, diags),
		Timeout:            flex.ExpandInt64Pointer(m.Timeout),
		VendorIdentifier:   flex.ExpandStringPointer(m.VendorIdentifier),
		WapiUserName:       flex.ExpandStringPointer(m.WapiUserName),
		WapiUserPassword:   flex.ExpandStringPointer(m.WapiUserPassword),
	}
	return to
}

func FlattenSyslogEndpoint(ctx context.Context, from *misc.SyslogEndpoint, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(SyslogEndpointAttrTypes)
	}
	m := SyslogEndpointModel{}
	m.Flatten(ctx, from, diags)
	m.ExtAttrsAll = types.MapNull(types.StringType)
	t, d := types.ObjectValueFrom(ctx, SyslogEndpointAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *SyslogEndpointModel) Flatten(ctx context.Context, from *misc.SyslogEndpoint, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = SyslogEndpointModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.LogLevel = flex.FlattenStringPointer(from.LogLevel)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.OutboundMemberType = flex.FlattenStringPointer(from.OutboundMemberType)
	m.OutboundMembers = flex.FlattenFrameworkListString(ctx, from.OutboundMembers, diags)
	planSyslogServers := m.SyslogServers
	m.SyslogServers = flex.FlattenFrameworkListNestedBlock(ctx, from.SyslogServers, SyslogEndpointSyslogServersAttrTypes, diags, FlattenSyslogEndpointSyslogServers)
	m.TemplateInstance = FlattenSyslogEndpointTemplateInstance(ctx, from.TemplateInstance, diags)
	m.Timeout = flex.FlattenInt64Pointer(from.Timeout)
	m.VendorIdentifier = flex.FlattenStringPointer(from.VendorIdentifier)
	m.WapiUserName = flex.FlattenStringPointer(from.WapiUserName)
	// Preserve WapiUserPassword - API doesn't return sensitive password fields
	// The value is already set from plan/state, so we don't overwrite it
	if !planSyslogServers.IsNull() {
		result, diags := utils.CopyFieldFromPlanToRespList(ctx, planSyslogServers, m.SyslogServers, "certificate_file_path")
		if !diags.HasError() {
			m.SyslogServers = result.(basetypes.ListValue)
		}
	}
}

func (m *SyslogEndpointModel) PutExpand(to *misc.SyslogEndpoint) *misc.SyslogEndpoint {
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

	for field, attr := range SyslogEndpointResourceSchemaAttributes {
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
