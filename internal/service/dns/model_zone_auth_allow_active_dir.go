package dns

import (
	"context"
	"fmt"
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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	customvalidator "github.com/infobloxopen/terraform-provider-nios/internal/validator"
)

type ZoneAuthAllowActiveDirModel struct {
	Address    types.String `tfsdk:"address"`
	Permission types.String `tfsdk:"permission"`
}

var ZoneAuthAllowActiveDirAttrTypes = map[string]attr.Type{
	"address":    types.StringType,
	"permission": types.StringType,
}

var ZoneAuthAllowActiveDirResourceSchemaAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The address this rule applies to or \"Any\".",
		Validators: []validator.String{
			customvalidator.ValidateTrimmedString(),
		},
	},
	"permission": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString("ALLOW"),
		Validators: []validator.String{
			stringvalidator.OneOf("ALLOW", "DENY"),
		},
		MarkdownDescription: "The permission to use for this address.",
	},
}

func ExpandZoneAuthAllowActiveDir(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneAuthAllowActiveDir {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneAuthAllowActiveDirModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneAuthAllowActiveDirModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneAuthAllowActiveDir {
	if m == nil {
		return nil
	}
	to := &dns.ZoneAuthAllowActiveDir{
		Address:    flex.ExpandStringPointer(m.Address),
		Permission: flex.ExpandStringPointer(m.Permission),
	}
	return to
}

func FlattenZoneAuthAllowActiveDir(ctx context.Context, from *dns.ZoneAuthAllowActiveDir, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneAuthAllowActiveDirAttrTypes)
	}
	m := ZoneAuthAllowActiveDirModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneAuthAllowActiveDirAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneAuthAllowActiveDirModel) Flatten(ctx context.Context, from *dns.ZoneAuthAllowActiveDir, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneAuthAllowActiveDirModel{}
	}
	m.Address = flex.FlattenStringPointer(from.Address)
	m.Permission = flex.FlattenStringPointer(from.Permission)
}

func (m *ZoneAuthAllowActiveDirModel) PutExpand(to *dns.ZoneAuthAllowActiveDir) *dns.ZoneAuthAllowActiveDir {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range ZoneAuthAllowActiveDirResourceSchemaAttributes {
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
							fmt.Printf("Field: %s, ok: %v, Computed: %v, fieldValue: %v, Value: %s\n", field, ok, boolComp, fieldValue, txtFieldValue)
							if ok {
								if boolComp && txtFieldValue == "" {
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
