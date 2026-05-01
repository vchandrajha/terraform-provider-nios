package cloud

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/cloud"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type AwsuserModel struct {
	Ref             types.String `tfsdk:"ref"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccountId       types.String `tfsdk:"account_id"`
	GovcloudEnabled types.Bool   `tfsdk:"govcloud_enabled"`
	LastUsed        types.Int64  `tfsdk:"last_used"`
	Name            types.String `tfsdk:"name"`
	NiosUserName    types.String `tfsdk:"nios_user_name"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	Status          types.String `tfsdk:"status"`
}

var AwsuserAttrTypes = map[string]attr.Type{
	"ref":               types.StringType,
	"access_key_id":     types.StringType,
	"account_id":        types.StringType,
	"govcloud_enabled":  types.BoolType,
	"last_used":         types.Int64Type,
	"name":              types.StringType,
	"nios_user_name":    types.StringType,
	"secret_access_key": types.StringType,
	"status":            types.StringType,
}

var AwsuserResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"access_key_id": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthAtMost(255),
		},
		MarkdownDescription: "The unique Access Key ID of this AWS user. Maximum 255 characters.",
	},
	"account_id": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtMost(64),
		},
		MarkdownDescription: "The AWS Account ID of this AWS user. Maximum 64 characters.",
	},
	"govcloud_enabled": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Indicates if gov cloud is enabled or disabled.",
	},
	"last_used": schema.Int64Attribute{
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The timestamp when this AWS user credentials was last used.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtMost(64),
		},
		MarkdownDescription: "The AWS user name. Maximum 64 characters.",
	},
	"nios_user_name": schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.LengthAtMost(64),
		},
		MarkdownDescription: "The NIOS user name mapped to this AWS user. Maximum 64 characters.",
	},
	"secret_access_key": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The Secret Access Key for the Access Key ID of this user. Maximum 255 characters.",
	},
	"status": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Indicate the validity status of this AWS user.",
	},
}

func (m *AwsuserModel) Expand(ctx context.Context, diags *diag.Diagnostics) *cloud.Awsuser {
	if m == nil {
		return nil
	}
	to := &cloud.Awsuser{
		AccessKeyId:     flex.ExpandStringPointer(m.AccessKeyId),
		AccountId:       flex.ExpandStringPointer(m.AccountId),
		GovcloudEnabled: flex.ExpandBoolPointer(m.GovcloudEnabled),
		Name:            flex.ExpandStringPointer(m.Name),
		NiosUserName:    flex.ExpandStringPointer(m.NiosUserName),
		SecretAccessKey: flex.ExpandStringPointer(m.SecretAccessKey),
	}
	return to
}

func FlattenAwsuser(ctx context.Context, from *cloud.Awsuser, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(AwsuserAttrTypes)
	}
	m := AwsuserModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, AwsuserAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *AwsuserModel) Flatten(ctx context.Context, from *cloud.Awsuser, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = AwsuserModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AccessKeyId = flex.FlattenStringPointer(from.AccessKeyId)
	m.AccountId = flex.FlattenStringPointer(from.AccountId)
	m.GovcloudEnabled = types.BoolPointerValue(from.GovcloudEnabled)
	m.LastUsed = flex.FlattenInt64Pointer(from.LastUsed)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.NiosUserName = flex.FlattenStringPointer(from.NiosUserName)
	m.Status = flex.FlattenStringPointer(from.Status)
}

func (m *AwsuserModel) PutExpand(to *cloud.Awsuser) *cloud.Awsuser {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range AwsuserResourceSchemaAttributes {
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
