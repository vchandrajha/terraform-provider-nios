package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/infobloxopen/infoblox-nios-go-client/dns"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type ZoneauthdnsseckeyparamsZskAlgorithmsModel struct {
	Algorithm types.String `tfsdk:"algorithm"`
	Size      types.Int64  `tfsdk:"size"`
}

var ZoneauthdnsseckeyparamsZskAlgorithmsAttrTypes = map[string]attr.Type{
	"algorithm": types.StringType,
	"size":      types.Int64Type,
}

var ZoneauthdnsseckeyparamsZskAlgorithmsResourceSchemaAttributes = map[string]schema.Attribute{
	"algorithm": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ECDSAP256SHA256", "ECDSAP384SHA384", "RSASHA1", "RSASHA256", "RSASHA512"),
		},
		MarkdownDescription: "The signing key algorithm.",
	},
	"size": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The signing key size, in bits.",
	},
}

func ExpandZoneauthdnsseckeyparamsZskAlgorithms(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.ZoneauthdnsseckeyparamsZskAlgorithms {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m ZoneauthdnsseckeyparamsZskAlgorithmsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *ZoneauthdnsseckeyparamsZskAlgorithmsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.ZoneauthdnsseckeyparamsZskAlgorithms {
	if m == nil {
		return nil
	}
	to := &dns.ZoneauthdnsseckeyparamsZskAlgorithms{
		Algorithm: flex.ExpandStringPointer(m.Algorithm),
		Size:      flex.ExpandInt64Pointer(m.Size),
	}
	return to
}

func FlattenZoneauthdnsseckeyparamsZskAlgorithms(ctx context.Context, from *dns.ZoneauthdnsseckeyparamsZskAlgorithms, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(ZoneauthdnsseckeyparamsZskAlgorithmsAttrTypes)
	}
	m := ZoneauthdnsseckeyparamsZskAlgorithmsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, ZoneauthdnsseckeyparamsZskAlgorithmsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *ZoneauthdnsseckeyparamsZskAlgorithmsModel) Flatten(ctx context.Context, from *dns.ZoneauthdnsseckeyparamsZskAlgorithms, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = ZoneauthdnsseckeyparamsZskAlgorithmsModel{}
	}
	m.Algorithm = flex.FlattenStringPointer(from.Algorithm)
	m.Size = flex.FlattenInt64Pointer(from.Size)
}

func (m *ZoneauthdnsseckeyparamsZskAlgorithmsModel) PutExpand(to *dns.ZoneauthdnsseckeyparamsZskAlgorithms) *dns.ZoneauthdnsseckeyparamsZskAlgorithms {
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

	for field, attr := range ZoneauthdnsseckeyparamsZskAlgorithmsResourceSchemaAttributes {
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
