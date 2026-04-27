package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ToString(field string, fieldValue any) string {
	if fieldValue == nil {
		return ""
	}
	var txtFieldValue string
	switch v := fieldValue.(type) {
	case *string:
		if v != nil {
			txtFieldValue = *v
		} else {
			txtFieldValue = ""
		}
	case string:
		txtFieldValue = v
	case *int:
		if v != nil {
			txtFieldValue = fmt.Sprintf("%d", *v)
		} else {
			txtFieldValue = ""
		}
	case int:
		txtFieldValue = fmt.Sprintf("%d", v)
	case *int64:
		if v != nil {
			txtFieldValue = fmt.Sprintf("%d", *v)
		} else {
			txtFieldValue = ""
		}
	case int64:
		txtFieldValue = fmt.Sprintf("%d", v)
	case *float64:
		if v != nil {
			txtFieldValue = fmt.Sprintf("%f", *v)
		} else {
			txtFieldValue = ""
		}
	case float64:
		txtFieldValue = fmt.Sprintf("%f", v)
	case *bool:
		if v != nil {
			txtFieldValue = fmt.Sprintf("%t", *v)
		} else {
			txtFieldValue = ""
		}
	case bool:
		txtFieldValue = fmt.Sprintf("%t", v)
	case nil:
		txtFieldValue = ""
	default:
		// Handle slices, arrays, maps, structs, pointers, etc.
		rv := reflect.ValueOf(fieldValue)
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				txtFieldValue = ""
			} else {
				elem := rv.Elem()
				if elem.Kind() == reflect.String {
					txtFieldValue = elem.String()
				} else if elem.Kind() == reflect.Int || elem.Kind() == reflect.Int64 {
					txtFieldValue = fmt.Sprintf("%d", elem.Int())
				} else if elem.Kind() == reflect.Float64 {
					txtFieldValue = fmt.Sprintf("%f", elem.Float())
				} else if elem.Kind() == reflect.Bool {
					txtFieldValue = fmt.Sprintf("%t", elem.Bool())
				} else {
					txtFieldValue = fmt.Sprintf("%v", elem.Interface())
				}
			}
		case reflect.Slice, reflect.Array:
			if rv.Len() == 0 {
				txtFieldValue = ""
			} else {
				txtFieldValue = fmt.Sprintf("%v", fieldValue)
			}
		case reflect.Map:
			if rv.Len() == 0 {
				txtFieldValue = ""
			} else {
				txtFieldValue = fmt.Sprintf("%v", fieldValue)
			}
		case reflect.Struct:
			txtFieldValue = fmt.Sprintf("%v", fieldValue)
		default:
			txtFieldValue = fmt.Sprintf("%v", fieldValue)
		}
	}
	return txtFieldValue
}

func DeleteBy(to any, field string) {
	// Use reflection to set the field to its zero value
	// Only works for exported fields with matching names (case-insensitive)
	val := reflect.ValueOf(to).Elem()
	for i := 0; i < val.NumField(); i++ {
		sf := val.Type().Field(i)
		if strings.EqualFold(sf.Name, field) {
			val.Field(i).Set(reflect.Zero(sf.Type))
			break
		}
	}
}
