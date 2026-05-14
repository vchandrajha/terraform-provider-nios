package dtc

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/dtc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	planmodifiers "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/immutable"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
	derivedmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/derived"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type DtcRecordCnameModel struct {
	Ref          types.String `tfsdk:"ref"`
	AutoCreated  types.Bool   `tfsdk:"auto_created"`
	Canonical    types.String `tfsdk:"canonical"`
	Comment      types.String `tfsdk:"comment"`
	Disable      types.Bool   `tfsdk:"disable"`
	DnsCanonical types.String `tfsdk:"dns_canonical"`
	DtcServer    types.String `tfsdk:"dtc_server"`
	Ttl          types.Int64  `tfsdk:"ttl"`
	UseTtl       types.Bool   `tfsdk:"use_ttl"`
}

var DtcRecordCnameAttrTypes = map[string]attr.Type{
	"ref":           types.StringType,
	"auto_created":  types.BoolType,
	"canonical":     types.StringType,
	"comment":       types.StringType,
	"disable":       types.BoolType,
	"dns_canonical": types.StringType,
	"dtc_server":    types.StringType,
	"ttl":           types.Int64Type,
	"use_ttl":       types.BoolType,
}

var DtcRecordCnameResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"auto_created": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Flag that indicates whether this record was automatically created by NIOS.",
	},
	"canonical": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The canonical name of the host.",
	},
	"comment": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 256),
			customvalidator.ValidateTrimmedString(),
		},
		MarkdownDescription: "Comment for the record; maximum 256 characters.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the record is disabled or not. False means that the record is enabled.",
	},
	"dns_canonical": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			derivedmod.PunycodeDerivedFrom("canonical"),
		},
		MarkdownDescription: "The canonical name as server by DNS protocol.",
	},
	"dtc_server": schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
		MarkdownDescription: "The name of the DTC Server object with which the DTC record is associated.",
	},
	"ttl": schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AlsoRequires(path.MatchRoot("use_ttl")),
		},
		MarkdownDescription: "The Time to Live (TTL) value.",
	},
	"use_ttl": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Use flag for: ttl",
	},
}

func (m *DtcRecordCnameModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *dtc.DtcRecordCname {
	if m == nil {
		return nil
	}
	to := &dtc.DtcRecordCname{
		Canonical: flex.ExpandStringPointer(m.Canonical),
		Comment:   flex.ExpandStringPointer(m.Comment),
		Disable:   flex.ExpandBoolPointer(m.Disable),
		Ttl:       flex.ExpandInt64Pointer(m.Ttl),
		UseTtl:    flex.ExpandBoolPointer(m.UseTtl),
	}
	if isCreate {
		to.DtcServer = flex.ExpandStringPointer(m.DtcServer)
	}
	return to
}

func FlattenDtcRecordCname(ctx context.Context, from *dtc.DtcRecordCname, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(DtcRecordCnameAttrTypes)
	}
	m := DtcRecordCnameModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, DtcRecordCnameAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *DtcRecordCnameModel) Flatten(ctx context.Context, from *dtc.DtcRecordCname, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = DtcRecordCnameModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AutoCreated = types.BoolPointerValue(from.AutoCreated)
	m.Canonical = flex.FlattenStringPointer(from.Canonical)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.DnsCanonical = flex.FlattenStringPointer(from.DnsCanonical)
	m.DtcServer = flex.FlattenStringPointer(from.DtcServer)
	m.Ttl = flex.FlattenInt64Pointer(from.Ttl)
	m.UseTtl = types.BoolPointerValue(from.UseTtl)
}

func (m *DtcRecordCnameModel) PutExpand(to *dtc.DtcRecordCname) *dtc.DtcRecordCname {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range DtcRecordCnameResourceSchemaAttributes {
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
