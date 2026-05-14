package grid

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/grid"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type MemberBgpAsModel struct {
	As         types.Int64 `tfsdk:"as"`
	Keepalive  types.Int64 `tfsdk:"keepalive"`
	Holddown   types.Int64 `tfsdk:"holddown"`
	Neighbors  types.List  `tfsdk:"neighbors"`
	LinkDetect types.Bool  `tfsdk:"link_detect"`
}

var MemberBgpAsAttrTypes = map[string]attr.Type{
	"as":          types.Int64Type,
	"keepalive":   types.Int64Type,
	"holddown":    types.Int64Type,
	"neighbors":   types.ListType{ElemType: types.ObjectType{AttrTypes: MemberbgpasNeighborsAttrTypes}},
	"link_detect": types.BoolType,
}

var MemberBgpAsResourceSchemaAttributes = map[string]schema.Attribute{
	"as": schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The number of this autonomous system.",
	},
	"keepalive": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(4),
		MarkdownDescription: "The AS keepalive timer (in seconds). The valid value is from 1 to 21845.",
	},
	"holddown": schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		Default:             int64default.StaticInt64(16),
		MarkdownDescription: "The AS holddown timer (in seconds). The valid value is from 3 to 65535.",
	},
	"neighbors": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: MemberbgpasNeighborsResourceSchemaAttributes,
		},
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The BGP neighbors for this AS.",
	},
	"link_detect": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if link detection on the interface is enabled or not.",
	},
}

func ExpandMemberBgpAs(ctx context.Context, o types.Object, diags *diag.Diagnostics) *grid.MemberBgpAs {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m MemberBgpAsModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *MemberBgpAsModel) Expand(ctx context.Context, diags *diag.Diagnostics) *grid.MemberBgpAs {
	if m == nil {
		return nil
	}
	to := &grid.MemberBgpAs{
		As:         flex.ExpandInt64Pointer(m.As),
		Keepalive:  flex.ExpandInt64Pointer(m.Keepalive),
		Holddown:   flex.ExpandInt64Pointer(m.Holddown),
		Neighbors:  flex.ExpandFrameworkListNestedBlock(ctx, m.Neighbors, diags, ExpandMemberbgpasNeighbors),
		LinkDetect: flex.ExpandBoolPointer(m.LinkDetect),
	}
	return to
}

func FlattenMemberBgpAs(ctx context.Context, from *grid.MemberBgpAs, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(MemberBgpAsAttrTypes)
	}
	m := MemberBgpAsModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, MemberBgpAsAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *MemberBgpAsModel) Flatten(ctx context.Context, from *grid.MemberBgpAs, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = MemberBgpAsModel{}
	}
	m.As = flex.FlattenInt64Pointer(from.As)
	m.Keepalive = flex.FlattenInt64Pointer(from.Keepalive)
	m.Holddown = flex.FlattenInt64Pointer(from.Holddown)
	m.Neighbors = flex.FlattenFrameworkListNestedBlock(ctx, from.Neighbors, MemberbgpasNeighborsAttrTypes, diags, FlattenMemberbgpasNeighbors)
	m.LinkDetect = types.BoolPointerValue(from.LinkDetect)
}

func (m *MemberBgpAsModel) PutExpand(to *grid.MemberBgpAs) *grid.MemberBgpAs {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range MemberBgpAsResourceSchemaAttributes {
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
