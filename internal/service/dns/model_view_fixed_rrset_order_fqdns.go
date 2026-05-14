package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type ViewFixedRrsetOrderFqdnsModel struct {
	Fqdn       types.String `tfsdk:"fqdn"`
	RecordType types.String `tfsdk:"record_type"`
}

var ViewFixedRrsetOrderFqdnsAttrTypes = map[string]attr.Type{
	"fqdn":        types.StringType,
	"record_type": types.StringType,
}

var ViewFixedRrsetOrderFqdnsResourceSchemaAttributes = map[string]schema.Attribute{
	"fqdn": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			customvalidator.IsValidDomainName(),
		},
		MarkdownDescription: "The FQDN of the fixed RRset configuration item.",
	},
	"record_type": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("A", "AAAA", "BOTH"),
		},
		Default:             stringdefault.StaticString("A"),
		MarkdownDescription: "The record type for the specified FQDN in the fixed RRset configuration.",
	},
}

func ExpandViewFixedRrsetOrderFqdns(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ViewFixedRrsetOrderFqdns {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ViewFixedRrsetOrderFqdnsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ViewFixedRrsetOrderFqdnsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ViewFixedRrsetOrderFqdns {
	if m == nil {
		return nil
	}
	to := &dns.ViewFixedRrsetOrderFqdns{
		Fqdn:       flex.ExpandStringPointer(m.Fqdn),
		RecordType: flex.ExpandStringPointer(m.RecordType),
	}
	return to
}

func FlattenViewFixedRrsetOrderFqdns(ctx context.Context, from *dns.ViewFixedRrsetOrderFqdns, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ViewFixedRrsetOrderFqdnsAttrTypes)
	}
	m := ViewFixedRrsetOrderFqdnsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ViewFixedRrsetOrderFqdnsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ViewFixedRrsetOrderFqdnsModel) Flatten(ctx context.Context, from *dns.ViewFixedRrsetOrderFqdns, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ViewFixedRrsetOrderFqdnsModel{}
	}
	m.Fqdn = flex.FlattenStringPointer(from.Fqdn)
	m.RecordType = flex.FlattenStringPointer(from.RecordType)
}

func (m *ViewFixedRrsetOrderFqdnsModel) PutExpand(to *dns.ViewFixedRrsetOrderFqdns) *dns.ViewFixedRrsetOrderFqdns {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ViewFixedRrsetOrderFqdnsResourceSchemaAttributes {
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
