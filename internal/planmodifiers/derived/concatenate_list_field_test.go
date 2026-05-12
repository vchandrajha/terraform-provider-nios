package derived

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExtractAndConcatenateFromList(t *testing.T) {
	elemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"field_type":     types.StringType,
			"field_value":    types.StringType,
			"include_length": types.StringType,
		},
	}

	obj1, _ := types.ObjectValue(elemType.AttrTypes, map[string]attr.Value{
		"field_type":     types.StringValue("T"),
		"field_value":    types.StringValue("hello"),
		"include_length": types.StringValue("8_BIT"),
	})
	obj2, _ := types.ObjectValue(elemType.AttrTypes, map[string]attr.Value{
		"field_type":     types.StringValue("T"),
		"field_value":    types.StringValue("world"),
		"include_length": types.StringValue("8_BIT"),
	})

	listVal, _ := types.ListValue(elemType, []attr.Value{obj1, obj2})

	tests := []struct {
		name      string
		key       string
		separator string
		expected  string
	}{
		{"concatenate field_value no sep", "field_value", "", "helloworld"},
		{"concatenate field_value with sep", "field_value", ",", "hello,world"},
		{"concatenate field_type", "field_type", "", "TT"},
		{"concatenate include_length", "include_length", "-", "8_BIT-8_BIT"},
		{"missing key", "nonexistent", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAndConcatenateFromList(listVal, tt.key, tt.separator)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractAndConcatenateFromList_NullList(t *testing.T) {
	nullList := types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{}})
	result := ExtractAndConcatenateFromList(nullList, "field_value", "")
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestExtractValuesFromList(t *testing.T) {
	elemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"field_value": types.StringType,
		},
	}

	obj1, _ := types.ObjectValue(elemType.AttrTypes, map[string]attr.Value{
		"field_value": types.StringValue("a"),
	})
	obj2, _ := types.ObjectValue(elemType.AttrTypes, map[string]attr.Value{
		"field_value": types.StringValue("b"),
	})

	listVal, _ := types.ListValue(elemType, []attr.Value{obj1, obj2})

	values := ExtractValuesFromList(listVal, "field_value")
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}

	for i, expected := range []string{"a", "b"} {
		strVal, ok := values[i].(types.String)
		if !ok {
			t.Fatalf("value %d is not types.String", i)
		}
		if strVal.ValueString() != expected {
			t.Errorf("value %d: got %q, want %q", i, strVal.ValueString(), expected)
		}
	}
}

func TestConcatenateStringListValues(t *testing.T) {
	listVal, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("hello"),
		types.StringValue("world"),
	})

	tests := []struct {
		name        string
		separator   string
		quoteValues bool
		expected    string
	}{
		{"no quotes no sep", "", false, "helloworld"},
		{"no quotes space sep", " ", false, "hello world"},
		{"quoted space sep", " ", true, `"hello" "world"`},
		{"quoted comma sep", ",", true, `"hello","world"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConcatenateStringListValues(listVal, tt.separator, tt.quoteValues)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConcatenateStringListValues_NullList(t *testing.T) {
	nullList := types.ListNull(types.StringType)
	result := ConcatenateStringListValues(nullList, " ", false)
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestConcatenateStringListValues_SingleElement(t *testing.T) {
	listVal, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("only"),
	})

	result := ConcatenateStringListValues(listVal, " ", true)
	if result != `"only"` {
		t.Errorf("got %q, want %q", result, `"only"`)
	}
}
