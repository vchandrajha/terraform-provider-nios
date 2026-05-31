package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberCspMemberSettingModel struct {
	UseCspJoinToken   types.Bool   `tfsdk:"use_csp_join_token"`
	UseCspDnsResolver types.Bool   `tfsdk:"use_csp_dns_resolver"`
	UseCspHttpsProxy  types.Bool   `tfsdk:"use_csp_https_proxy"`
	CspJoinToken      types.String `tfsdk:"csp_join_token"`
	CspDnsResolver    types.String `tfsdk:"csp_dns_resolver"`
	CspHttpsProxy     types.String `tfsdk:"csp_https_proxy"`
}

var MemberCspMemberSettingAttrTypes = map[string]attr.Type{
	"use_csp_join_token":   types.BoolType,
	"use_csp_dns_resolver": types.BoolType,
	"use_csp_https_proxy":  types.BoolType,
	"csp_join_token":       types.StringType,
	"csp_dns_resolver":     types.StringType,
	"csp_https_proxy":      types.StringType,
}

var MemberCspMemberSettingResourceSchemaAttributes = map[string]schema.Attribute{
	"use_csp_join_token": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Overrides grid join token",
	},
	"use_csp_dns_resolver": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Overrides CSP DNS Resolver",
	},
	"use_csp_https_proxy": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Overrides grid https proxy",
	},
	"csp_join_token": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			stringvalidator.AlsoRequires(path.MatchRoot("use_csp_join_token")),
		},
		MarkdownDescription: "Join token required to connect to a cluster",
	},
	"csp_dns_resolver": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("use_csp_dns_resolver")),
		},
		MarkdownDescription: "IP address of DNS resolver in DFP",
	},
	"csp_https_proxy": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			stringvalidator.AlsoRequires(path.MatchRoot("use_csp_https_proxy")),
		},
		MarkdownDescription: "HTTP Proxy IP address of CSP Portal",
	},
}

func ExpandMemberCspMemberSetting(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberCspMemberSetting {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberCspMemberSettingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberCspMemberSettingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberCspMemberSetting {
	if m == nil {
		return nil
	}
	to := &grid.MemberCspMemberSetting{
		UseCspJoinToken:   flex.ExpandBoolPointer(m.UseCspJoinToken),
		UseCspDnsResolver: flex.ExpandBoolPointer(m.UseCspDnsResolver),
		UseCspHttpsProxy:  flex.ExpandBoolPointer(m.UseCspHttpsProxy),
		CspJoinToken:      flex.ExpandStringPointer(m.CspJoinToken),
		CspDnsResolver:    flex.ExpandStringPointerEmptyAsNil(m.CspDnsResolver),
		CspHttpsProxy:     flex.ExpandStringPointer(m.CspHttpsProxy),
	}
	return to
}

func FlattenMemberCspMemberSetting(ctx context.Context, from *grid.MemberCspMemberSetting, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberCspMemberSettingAttrTypes)
	}
	m := MemberCspMemberSettingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberCspMemberSettingAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberCspMemberSettingModel) Flatten(ctx context.Context, from *grid.MemberCspMemberSetting, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberCspMemberSettingModel{}
	}
	m.UseCspJoinToken = types.BoolPointerValue(from.UseCspJoinToken)
	m.UseCspDnsResolver = types.BoolPointerValue(from.UseCspDnsResolver)
	m.UseCspHttpsProxy = types.BoolPointerValue(from.UseCspHttpsProxy)
	m.CspJoinToken = flex.FlattenStringPointer(from.CspJoinToken)
	m.CspDnsResolver = flex.FlattenStringPointer(from.CspDnsResolver)
	m.CspHttpsProxy = flex.FlattenStringPointer(from.CspHttpsProxy)
}

func (m *MemberCspMemberSettingModel) PutExpand(to *grid.MemberCspMemberSetting) *grid.MemberCspMemberSetting {
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

	for field, attr := range MemberCspMemberSettingResourceSchemaAttributes {
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
