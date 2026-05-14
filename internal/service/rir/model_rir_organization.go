package rir

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/rir"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type RirOrganizationModel struct {
	Ref         types.String `tfsdk:"ref"`
	ExtAttrs    types.Map    `tfsdk:"extattrs"`
	Id          types.String `tfsdk:"id"`
	Maintainer  types.String `tfsdk:"maintainer"`
	Name        types.String `tfsdk:"name"`
	Password    types.String `tfsdk:"password"`
	Rir         types.String `tfsdk:"rir"`
	SenderEmail types.String `tfsdk:"sender_email"`
}

var RirOrganizationAttrTypes = map[string]attr.Type{
	"ref":          types.StringType,
	"extattrs":     types.MapType{ElemType: types.StringType},
	"id":           types.StringType,
	"maintainer":   types.StringType,
	"name":         types.StringType,
	"password":     types.StringType,
	"rir":          types.StringType,
	"sender_email": types.StringType,
}

var RirOrganizationResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"extattrs": schema.MapAttribute{
		Required:    true,
		ElementType: types.StringType,
		Validators: []validator.Map{
			customvalidator.MapContainsKey("RIPE Admin Contact"),
			customvalidator.MapContainsKey("RIPE Country"),
			customvalidator.MapContainsKey("RIPE Technical Contact"),
			customvalidator.MapContainsKey("RIPE Email"),
			mapvalidator.KeysAre(stringvalidator.OneOf("RIPE Description", "RIPE Admin Contact", "RIPE Country", "RIPE Technical Contact", "RIPE Email", "RIPE Remarks", "RIPE Notify", "RIPE Registry Source", "RIPE Organization Type", "RIPE Address", "RIPE Phone Number", "RIPE Fax Number", "RIPE Abuse Mailbox", "RIPE Reference Notify")),
		},
		MarkdownDescription: "Extensible attributes associated with the object.",
	},
	"id": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(`^ORG-[A-Za-z]{2,4}[1-9][0-9]{0,4}-[A-Za-z0-9]{1,9}$`), "- Invalid Organization ID. A Valid Organization ID starts with 'ORG-', followed by 2-4 letters, then a number between 1 and 99999, and ends with a hyphen and 1-9 alphanumeric characters. Valid Examples for ID are ORG-CA1-RIPE or ORG-CB2-TEST"),
		},
		MarkdownDescription: "The RIR organization identifier. Valid Examples for ID are ORG-CA1-RIPE or ORG-CB2-TEST ",
	},
	"maintainer": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 80),
			stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9]$`), "- A valid maintainer starts with a letter, followed by letters, numbers, underscores, or hyphens, and ends with a letter or number. Valid examples for maintainer are 'infoblox' and 'nios-support'"),
		},
		MarkdownDescription: "The RIR organization maintainer.",
	},
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
			stringvalidator.LengthBetween(0, 256),
			stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "- Invalid Organization Name. A valid organization name can only contain letters, numbers, underscores, or hyphens."),
		},
		MarkdownDescription: "The RIR organization name.",
	},
	"password": schema.StringAttribute{
		Required:  true,
		Sensitive: true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "The password for the maintainer of RIR organization.",
	},
	"rir": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("RIPE"),
		Validators: []validator.String{
			stringvalidator.OneOf("RIPE"),
		},
		MarkdownDescription: "The RIR associated with RIR organization.",
	},
	"sender_email": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(`^[^@]+@[^@]+\.com$`), "- must be a valid .com email address"),
		},
		MarkdownDescription: "The sender e-mail address for RIR organization.",
	},
}

func (m *RirOrganizationModel) Expand(ctx context.Context, diags *diag.Diagnostics) *rir.RirOrganization {
	if m == nil {
		return nil
	}
	to := &rir.RirOrganization{
		Id:          flex.ExpandStringPointer(m.Id),
		Maintainer:  flex.ExpandStringPointer(m.Maintainer),
		ExtAttrs:    ExpandExtAttrs(ctx, m.ExtAttrs, diags),
		Name:        flex.ExpandStringPointer(m.Name),
		Password:    flex.ExpandStringPointer(m.Password),
		Rir:         flex.ExpandStringPointer(m.Rir),
		SenderEmail: flex.ExpandStringPointer(m.SenderEmail),
	}
	return to
}

func FlattenRirOrganization(ctx context.Context, from *rir.RirOrganization, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RirOrganizationAttrTypes)
	}
	m := RirOrganizationModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RirOrganizationAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RirOrganizationModel) Flatten(ctx context.Context, from *rir.RirOrganization, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RirOrganizationModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.ExtAttrs = FlattenExtAttrs(ctx, m.ExtAttrs, from.ExtAttrs, diags)
	m.Id = flex.FlattenStringPointer(from.Id)
	m.Maintainer = flex.FlattenStringPointer(from.Maintainer)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Rir = flex.FlattenStringPointer(from.Rir)
	m.SenderEmail = flex.FlattenStringPointer(from.SenderEmail)
}

func (m *RirOrganizationModel) PutExpand(to *rir.RirOrganization) *rir.RirOrganization {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RirOrganizationResourceSchemaAttributes {
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
