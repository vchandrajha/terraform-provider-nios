package misc

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infobloxopen/infoblox-nios-go-client/misc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	planmodifiers "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/immutable"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
	refmod "github.com/infobloxopen/terraform-provider-nios/internal/planmodifiers/ref"
)

type TftpfiledirModel struct {
	Ref             types.String `tfsdk:"ref"`
	Directory       types.String `tfsdk:"directory"`
	IsSyncedToGm    types.Bool   `tfsdk:"is_synced_to_gm"`
	LastModify      types.Int64  `tfsdk:"last_modify"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	VtftpDirMembers types.List   `tfsdk:"vtftp_dir_members"`
}

var TftpfiledirAttrTypes = map[string]attr.Type{
	"ref":               types.StringType,
	"directory":         types.StringType,
	"is_synced_to_gm":   types.BoolType,
	"last_modify":       types.Int64Type,
	"name":              types.StringType,
	"type":              types.StringType,
	"vtftp_dir_members": types.ListType{ElemType: types.ObjectType{AttrTypes: TftpfiledirVtftpDirMembersAttrTypes}},
}

var TftpfiledirResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			refmod.UseStateUnlessResourceChanges(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"directory": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("/"),
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
		MarkdownDescription: "The path to the directory that contains file or subdirectory.",
	},
	"is_synced_to_gm": schema.BoolAttribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "Determines whether the TFTP entity is synchronized to Grid Master.",
	},
	"last_modify": schema.Int64Attribute{
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The time when the file or directory was last modified.",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The TFTP directory or file name.",
	},
	"type": schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			planmodifiers.ImmutableString(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("DIRECTORY", "FILE"),
		},
		MarkdownDescription: "The type of TFTP file system entity (directory or file). TYPE `FILE` is not supported through terraform provider and is reserved for future use.",
	},
	"vtftp_dir_members": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: TftpfiledirVtftpDirMembersResourceSchemaAttributes,
		},
		Computed: true,
		Optional: true,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		MarkdownDescription: "The replication members with TFTP client addresses where this virtual folder is applicable.",
	},
}

func (m *TftpfiledirModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *misc.Tftpfiledir {
	if m == nil {
		return nil
	}
	to := &misc.Tftpfiledir{
		Name: flex.ExpandStringPointer(m.Name),
	}
	if isCreate {
		to.Directory = flex.ExpandStringPointer(m.Directory)
		to.Type = flex.ExpandStringPointer(m.Type)
	}
	if !isCreate {
		to.VtftpDirMembers = flex.ExpandFrameworkListNestedBlock(ctx, m.VtftpDirMembers, diags, ExpandTftpfiledirVtftpDirMembers)
	}
	return to
}

func FlattenTftpfiledir(ctx context.Context, from *misc.Tftpfiledir, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(TftpfiledirAttrTypes)
	}
	m := TftpfiledirModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, TftpfiledirAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *TftpfiledirModel) Flatten(ctx context.Context, from *misc.Tftpfiledir, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = TftpfiledirModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.Directory = flex.FlattenStringPointer(from.Directory)
	m.IsSyncedToGm = types.BoolPointerValue(from.IsSyncedToGm)
	m.LastModify = flex.FlattenInt64Pointer(from.LastModify)
	m.Name = flex.FlattenStringPointer(from.Name)
	m.Type = flex.FlattenStringPointer(from.Type)
	m.VtftpDirMembers = flex.FlattenFrameworkListNestedBlock(ctx, from.VtftpDirMembers, TftpfiledirVtftpDirMembersAttrTypes, diags, FlattenTftpfiledirVtftpDirMembers)
}

func (m *TftpfiledirModel) PutExpand(to *misc.Tftpfiledir) *misc.Tftpfiledir {
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

	for field, attr := range TftpfiledirResourceSchemaAttributes {
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
