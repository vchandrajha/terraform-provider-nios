package dns

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/infoblox-nios-go-client/dns"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/infobloxopen/terraform-provider-nios/internal/flex"
	"github.com/infobloxopen/terraform-provider-nios/internal/utils"
)

type RecordHostSnmpCredentialModel struct {
	CommunityString types.String `tfsdk:"community_string"`
	Comment         types.String `tfsdk:"comment"`
	CredentialGroup types.String `tfsdk:"credential_group"`
}

var RecordHostSnmpCredentialAttrTypes = map[string]attr.Type{
	"community_string": types.StringType,
	"comment":          types.StringType,
	"credential_group": types.StringType,
}

var RecordHostSnmpCredentialResourceSchemaAttributes = map[string]schema.Attribute{
	"community_string": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The public community string.",
	},
	"comment": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Comments for the SNMPv1 and SNMPv2 users.",
	},
	"credential_group": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Group for the SNMPv1 and SNMPv2 credential.",
	},
}

func ExpandRecordHostSnmpCredential(ctx context.Context, o types.Object, diags *diag.Diagnostics) *dns.RecordHostSnmpCredential {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m RecordHostSnmpCredentialModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

func (m *RecordHostSnmpCredentialModel) Expand(ctx context.Context, diags *diag.Diagnostics) *dns.RecordHostSnmpCredential {
	if m == nil {
		return nil
	}
	to := &dns.RecordHostSnmpCredential{
		CommunityString: flex.ExpandStringPointer(m.CommunityString),
		Comment:         flex.ExpandStringPointer(m.Comment),
		CredentialGroup: flex.ExpandStringPointer(m.CredentialGroup),
	}
	return to
}

func FlattenRecordHostSnmpCredential(ctx context.Context, from *dns.RecordHostSnmpCredential, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RecordHostSnmpCredentialAttrTypes)
	}
	m := RecordHostSnmpCredentialModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RecordHostSnmpCredentialAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RecordHostSnmpCredentialModel) Flatten(ctx context.Context, from *dns.RecordHostSnmpCredential, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RecordHostSnmpCredentialModel{}
	}
	m.CommunityString = flex.FlattenStringPointer(from.CommunityString)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.CredentialGroup = flex.FlattenStringPointer(from.CredentialGroup)
}

func (m *RecordHostSnmpCredentialModel) PutExpand(to *dns.RecordHostSnmpCredential) *dns.RecordHostSnmpCredential {
	if m == nil {
		return nil
	}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range RecordHostSnmpCredentialResourceSchemaAttributes {
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
