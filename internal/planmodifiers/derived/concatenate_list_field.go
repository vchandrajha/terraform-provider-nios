package derived

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ planmodifier.String = concatenateListFieldModifier{}

// concatenateListFieldModifier computes a planned string value by extracting
// a specific key from each element of a list attribute and concatenating
// the values together.
//
// For example, given:
//
//	subfield_values = [
//	  { field_type = "T", field_value = "hello", include_length = "8_BIT" },
//	  { field_type = "T", field_value = "world", include_length = "8_BIT" },
//	]
//
// With listAttribute="subfield_values" and extractKey="field_value",
// the result would be "helloworld".
type concatenateListFieldModifier struct {
	listAttribute string
	extractKey    string
	separator     string
	quoteValues   bool
}

func (m concatenateListFieldModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Concatenates %q values from list attribute %q", m.extractKey, m.listAttribute)
}

func (m concatenateListFieldModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m concatenateListFieldModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the resource is being destroyed, do nothing.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Read the list attribute's planned value.
	var listVal types.List
	diags := req.Plan.GetAttribute(ctx, req.Path.ParentPath().AtName(m.listAttribute), &listVal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the list is unknown or null, leave the value unknown.
	if listVal.IsUnknown() || listVal.IsNull() {
		return
	}

	// Extract the key values from each list element and concatenate.
	var values []string
	for _, elem := range listVal.Elements() {
		objVal, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}

		attrs := objVal.Attributes()
		if val, exists := attrs[m.extractKey]; exists {
			if strVal, ok := val.(basetypes.StringValue); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				v := strVal.ValueString()
				if m.quoteValues {
					v = `"` + v + `"`
				}
				values = append(values, v)
			}
		}
	}

	resp.PlanValue = types.StringValue(strings.Join(values, m.separator))
}

// ConcatenateListField returns a plan modifier that computes a string value
// by extracting a specific key from each element of a list attribute and
// concatenating the values.
//
// Usage in schema:
//
//	"computed_field": schema.StringAttribute{
//	    Computed: true,
//	    PlanModifiers: []planmodifier.String{
//	        derivedmod.ConcatenateListField("subfield_values", "field_value", ""),
//	    },
//	}
//
// Parameters:
//   - listAttribute: the name of the sibling list attribute to read from
//   - extractKey: the key to extract from each object in the list
//   - separator: the string to join values with (use "" for no separator)
//   - quoteValues: if true, each value is wrapped in double quotes (e.g. `"val1" "val2"`)
func ConcatenateListField(listAttribute, extractKey, separator string, quoteValues bool) planmodifier.String {
	return concatenateListFieldModifier{
		listAttribute: listAttribute,
		extractKey:    extractKey,
		separator:     separator,
		quoteValues:   quoteValues,
	}
}
