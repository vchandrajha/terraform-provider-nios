package misc

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/misc"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type TftpfiledirVtftpDirMembersModel struct {
	Member       types.String      `tfsdk:"member"`
	IpType       types.String      `tfsdk:"ip_type"`
	Address      iptypes.IPAddress `tfsdk:"address"`
	StartAddress iptypes.IPAddress `tfsdk:"start_address"`
	EndAddress   iptypes.IPAddress `tfsdk:"end_address"`
	Network      types.String      `tfsdk:"network"`
	Cidr         types.Int64       `tfsdk:"cidr"`
}

var TftpfiledirVtftpDirMembersAttrTypes = map[string]attr.Type{
	"member":        types.StringType,
	"ip_type":       types.StringType,
	"address":       iptypes.IPAddressType{},
	"start_address": iptypes.IPAddressType{},
	"end_address":   iptypes.IPAddressType{},
	"network":       types.StringType,
	"cidr":          types.Int64Type,
}

var TftpfiledirVtftpDirMembersResourceSchemaAttributes = map[string]schema.Attribute{
	"member": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The Grid member on which to create the virtual TFTP directory.",
	},
	"ip_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("ADDRESS", "NETWORK", "RANGE"),
		},
		MarkdownDescription: "The IP type of the virtual TFTP root directory.",
	},
	"address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The IP address of the clients which will see the virtual TFTP directory as the root directory.",
	},
	"start_address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The start IP address of the range within which the clients will see the virtual TFTP directory as the root directory.",
	},
	"end_address": schema.StringAttribute{
		CustomType:          iptypes.IPAddressType{},
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The end IP address of the range within which the clients will see the virtual TFTP directory as the root directory.",
	},
	"network": schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "The IP address of network the clients from which will see the virtual TFTP directory as the root directory.",
	},
	"cidr": schema.Int64Attribute{
		Optional: true,
		Validators: []validator.Int64{
			int64validator.Between(0, 128),
		},
		MarkdownDescription: "The CIDR of network the clients from which will see the virtual TFTP directory as the root directory.",
	},
}

func ExpandTftpfiledirVtftpDirMembers(ctx context.Context, o types.Object, diags *diag.Diagnostics) *misc.TftpfiledirVtftpDirMembers {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m TftpfiledirVtftpDirMembersModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *TftpfiledirVtftpDirMembersModel) Expand(ctx context.Context, diags *diag.Diagnostics) *misc.TftpfiledirVtftpDirMembers {
	if m == nil {
		return nil
	}
	to := &misc.TftpfiledirVtftpDirMembers{
		Member:       flex.ExpandStringPointer(m.Member),
		IpType:       flex.ExpandStringPointer(m.IpType),
		Address:      flex.ExpandIPAddress(m.Address),
		StartAddress: flex.ExpandIPAddress(m.StartAddress),
		EndAddress:   flex.ExpandIPAddress(m.EndAddress),
		Network:      flex.ExpandStringPointer(m.Network),
		Cidr:         flex.ExpandInt64Pointer(m.Cidr),
	}
	return to
}

func FlattenTftpfiledirVtftpDirMembers(ctx context.Context, from *misc.TftpfiledirVtftpDirMembers, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(TftpfiledirVtftpDirMembersAttrTypes)
	}
	m := TftpfiledirVtftpDirMembersModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, TftpfiledirVtftpDirMembersAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *TftpfiledirVtftpDirMembersModel) Flatten(ctx context.Context, from *misc.TftpfiledirVtftpDirMembers, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = TftpfiledirVtftpDirMembersModel{}
	}
	m.Member = flex.FlattenStringPointer(from.Member)
	m.IpType = flex.FlattenStringPointer(from.IpType)
	m.Address = flex.FlattenIPAddress(from.Address)
	m.StartAddress = flex.FlattenIPAddress(from.StartAddress)
	m.EndAddress = flex.FlattenIPAddress(from.EndAddress)
	m.Network = flex.FlattenStringPointer(from.Network)
	m.Cidr = flex.FlattenInt64Pointer(from.Cidr)
}

func (m *TftpfiledirVtftpDirMembersModel) PutExpand(to *misc.TftpfiledirVtftpDirMembers) *misc.TftpfiledirVtftpDirMembers {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range TftpfiledirVtftpDirMembersResourceSchemaAttributes {
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
