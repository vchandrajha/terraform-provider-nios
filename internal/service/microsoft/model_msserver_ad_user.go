package microsoft

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/microsoft"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MsserverAdUserModel struct {
	LoginName                  types.String `tfsdk:"login_name"`
	LoginPassword              types.String `tfsdk:"login_password"`
	EnableUserSync             types.Bool   `tfsdk:"enable_user_sync"`
	SynchronizationInterval    types.Int64  `tfsdk:"synchronization_interval"`
	LastSyncTime               types.Int64  `tfsdk:"last_sync_time"`
	LastSyncStatus             types.String `tfsdk:"last_sync_status"`
	LastSyncDetail             types.String `tfsdk:"last_sync_detail"`
	LastSuccessSyncTime        types.Int64  `tfsdk:"last_success_sync_time"`
	UseLogin                   types.Bool   `tfsdk:"use_login"`
	UseEnableAdUserSync        types.Bool   `tfsdk:"use_enable_ad_user_sync"`
	UseSynchronizationMinDelay types.Bool   `tfsdk:"use_synchronization_min_delay"`
	UseEnableUserSync          types.Bool   `tfsdk:"use_enable_user_sync"`
	UseSynchronizationInterval types.Bool   `tfsdk:"use_synchronization_interval"`
}

var MsserverAdUserAttrTypes = map[string]attr.Type{
	"login_name":                    types.StringType,
	"login_password":                types.StringType,
	"enable_user_sync":              types.BoolType,
	"synchronization_interval":      types.Int64Type,
	"last_sync_time":                types.Int64Type,
	"last_sync_status":              types.StringType,
	"last_sync_detail":              types.StringType,
	"last_success_sync_time":        types.Int64Type,
	"use_login":                     types.BoolType,
	"use_enable_ad_user_sync":       types.BoolType,
	"use_synchronization_min_delay": types.BoolType,
	"use_enable_user_sync":          types.BoolType,
	"use_synchronization_interval":  types.BoolType,
}

var MsserverAdUserResourceSchemaAttributes = map[string]schema.Attribute{
	"login_name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The login name of the Microsoft Server.",
	},
	"login_password": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Sensitive:           true,
		MarkdownDescription: "The login password of the DHCP Microsoft Server.",
	},
	"enable_user_sync": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default:  booldefault.StaticBool(false),
		Validators: []validator.Bool{
			boolvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("use_enable_user_sync")),
		},
		MarkdownDescription: "Determines whether the Active Directory user synchronization is enabled or not.",
	},
	"synchronization_interval": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("use_synchronization_interval")),
		},
		MarkdownDescription: "The minimum number of minutes between two synchronizations.",
	},
	"last_sync_time": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Timestamp of the last synchronization attempt.",
	},
	"last_sync_status": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The status of the last synchronization attempt.",
	},
	"last_sync_detail": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The detailed status of the last synchronization attempt.",
	},
	"last_success_sync_time": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Timestamp of the last successful synchronization attempt.",
	},
	"use_login": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Flag to override login name and password from MS server",
	},
	"use_enable_ad_user_sync": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Flag to override AD User sync from grid level",
	},
	"use_synchronization_min_delay": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag to override synchronization interval from the MS Server",
	},
	"use_enable_user_sync": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Use flag for: enable_user_sync",
	},
	"use_synchronization_interval": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Use flag for: synchronization_interval",
	},
}

func ExpandMsserverAdUser(ctx context.Context, o types.Object, diags *diag.Diagnostics) *microsoft.MsserverAdUser {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MsserverAdUserModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MsserverAdUserModel) Expand(ctx context.Context, diags *diag.Diagnostics) *microsoft.MsserverAdUser {
	if m == nil {
		return nil
	}
	to := &microsoft.MsserverAdUser{
		LoginName:                  flex.ExpandStringPointer(m.LoginName),
		LoginPassword:              flex.ExpandStringPointer(m.LoginPassword),
		EnableUserSync:             flex.ExpandBoolPointer(m.EnableUserSync),
		SynchronizationInterval:    flex.ExpandInt64Pointer(m.SynchronizationInterval),
		UseLogin:                   flex.ExpandBoolPointer(m.UseLogin),
		UseEnableAdUserSync:        flex.ExpandBoolPointer(m.UseEnableAdUserSync),
		UseSynchronizationMinDelay: flex.ExpandBoolPointer(m.UseSynchronizationMinDelay),
		UseEnableUserSync:          flex.ExpandBoolPointer(m.UseEnableUserSync),
		UseSynchronizationInterval: flex.ExpandBoolPointer(m.UseSynchronizationInterval),
	}
	return to
}

func FlattenMsserverAdUser(ctx context.Context, from *microsoft.MsserverAdUser, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MsserverAdUserAttrTypes)
	}
	m := MsserverAdUserModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MsserverAdUserAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MsserverAdUserModel) Flatten(ctx context.Context, from *microsoft.MsserverAdUser, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MsserverAdUserModel{}
	}
	m.LoginName = flex.FlattenStringPointer(from.LoginName)
	m.LoginPassword = flex.FlattenStringPointer(from.LoginPassword)
	m.EnableUserSync = types.BoolPointerValue(from.EnableUserSync)
	m.SynchronizationInterval = flex.FlattenInt64Pointer(from.SynchronizationInterval)
	m.LastSyncTime = flex.FlattenInt64Pointer(from.LastSyncTime)
	m.LastSyncStatus = flex.FlattenStringPointer(from.LastSyncStatus)
	m.LastSyncDetail = flex.FlattenStringPointer(from.LastSyncDetail)
	m.LastSuccessSyncTime = flex.FlattenInt64Pointer(from.LastSuccessSyncTime)
	m.UseLogin = types.BoolPointerValue(from.UseLogin)
	m.UseEnableAdUserSync = types.BoolPointerValue(from.UseEnableAdUserSync)
	m.UseSynchronizationMinDelay = types.BoolPointerValue(from.UseSynchronizationMinDelay)
	m.UseEnableUserSync = types.BoolPointerValue(from.UseEnableUserSync)
	m.UseSynchronizationInterval = types.BoolPointerValue(from.UseSynchronizationInterval)
}

func (m *MsserverAdUserModel) PutExpand(to *microsoft.MsserverAdUser) *microsoft.MsserverAdUser {
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

	for field, attr := range MsserverAdUserResourceSchemaAttributes {
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
