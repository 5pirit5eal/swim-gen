package models

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// GetRequestJSON decodes a JSON request body into the provided struct.
// It also checks the Content-Type header to ensure it is "application/json".
// If the Content-Type is not "application/json", it returns an error.
func GetRequestJSON(r *http.Request, v any) error {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		return fmt.Errorf("unsupported content type: %s", mediaType)
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()

	return decoder.Decode(v)
}

// WriteResponseJSON writes a JSON response to the http.ResponseWriter.
// It sets the Content-Type header to "application/json" and writes the
// provided status code and value to the response body.
func WriteResponseJSON(w http.ResponseWriter, statusCode int, v any) error {
	json, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
	return nil
}

// StructToMap converts any struct to a map[string]any using the JSON tags
func StructToMap(item any) map[string]any {
	out := make(map[string]any)

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil
	}

	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		// gets us a struct field
		field := typ.Field(i)
		// reads the tag value
		tag := field.Tag.Get("json")
		// gets us a value
		value := v.Field(i).Interface()

		out[tag] = value
	}
	return out
}

func JSONInterfaceToStruct(data any, target any) error {
	dv := reflect.ValueOf(data)
	if dv.Kind() == reflect.Ptr {
		dv = dv.Elem()
	}
	tv := reflect.ValueOf(target)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}

	tt := tv.Type()

	switch {
	case dv.Kind() == reflect.Map && tt.Kind() == reflect.Struct:
		for i := 0; i < tt.NumField(); i++ {
			field := tt.Field(i)
			jsonTag, found := field.Tag.Lookup("json")
			if !found {
				return fmt.Errorf("error getting JSON tag")
			}
			fieldValue := dv.MapIndex(reflect.ValueOf(jsonTag))
			if fieldValue.IsValid() {
				if err := JSONInterfaceToStruct(fieldValue.Interface(), tv.FieldByName(field.Name).Addr().Interface()); err != nil {
					return err
				}
			}
		}
	case dv.Kind() == reflect.Slice && tt.Kind() == reflect.Slice:
		for i := 0; i < dv.Len(); i++ {
			elem := dv.Index(i)
			newElem := reflect.New(tt.Elem()).Elem()
			if err := JSONInterfaceToStruct(elem.Interface(), newElem.Addr().Interface()); err != nil {
				return err
			}
			tv.Set(reflect.Append(tv, newElem))
		}
	case dv.Kind() == reflect.Float64 && tt.Kind() == reflect.Int:
		tv.Set(reflect.ValueOf(int(dv.Float())))
	case dv.Kind() == tt.Kind():
		tv.Set(reflect.ValueOf(dv.Interface()))
	}
	return nil

}

// UnmarshalMapToStruct populates the fields of the struct pointed to by 'v'
// based on the key-value pairs in the 'data' map.
//
// It uses reflection to find struct fields. It first checks for a struct tag
// with the key 'csv' (e.g., `csv:"custom_key_name"`) to match
// the map key. If no 'csv' tag is found or it doesn't match, it falls
// back to matching the map key against the struct field name (case-sensitive).
//
// Supported field type conversions: string, int, float, bool.
//
// The 'v' argument must be a pointer to a struct.
func UnmarshalMapToStruct(data map[string]any, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v must be a non-nil pointer to a struct")
	}
	structVal := rv.Elem()
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}
	structType := structVal.Type()

	for i := range structVal.NumField() {
		fieldVal := structVal.Field(i)
		fieldType := structType.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		tag := fieldType.Tag.Get("csv")
		tag = strings.Split(tag, ",")[0] // Get the first part of the tag
		fieldName := fieldType.Name

		fieldFound := false
		for key, value := range data {

			match := false
			if tag != "" && tag == key {
				match = true
			} else if tag == "" && strings.EqualFold(fieldName, key) {
				match = true
			}
			if !match {
				continue
			}
			fieldFound = true
			err := setField(fieldVal, reflect.ValueOf(value))
			if err != nil {
				return fmt.Errorf("error setting field %s: %w", fieldName, err)
			}
			break
		}
		if !fieldFound {
			if !strings.Contains(fieldType.Tag.Get("csv"), "omitempty") {
				return fmt.Errorf("field %s mandatory in struct, not found in %v", fieldName, data)
			}
		}
	}
	return nil
}

// setField attempts to convert the type of 'value' to the type of 'fieldVal'
// and set the field's value.
func setField(fieldVal reflect.Value, value reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panic gracefully
			err = fmt.Errorf("panic occurred while setting field: %v", r)
		}
	}()

	fieldKind := fieldVal.Kind()

	switch fieldKind {
	case value.Kind():
		fieldVal.Set(value)
	case reflect.Int:
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue := value.Int()
			fieldVal.Set(reflect.ValueOf(int(intValue)))
		case reflect.Float32, reflect.Float64:
			floatValue := value.Float()
			fieldVal.Set(reflect.ValueOf(int(floatValue)))
		case reflect.String:
			strValue := value.String()
			intValue, err := strconv.Atoi(strValue)
			if err != nil {
				return fmt.Errorf("error converting string to int: %w", err)
			}
			fieldVal.Set(reflect.ValueOf(intValue))
		default:
			return fmt.Errorf("unsupported type conversion: %s to %s", value.Kind(), fieldKind)
		}
	case reflect.String:
		fieldVal.Set(value.Convert(reflect.TypeOf("")))
	case reflect.Bool:
		switch value.Kind() {
		case reflect.String:
			strValue := value.String()
			boolValue, err := strconv.ParseBool(strValue)
			if err != nil {
				return fmt.Errorf("error converting string to bool: %w", err)
			}
			fieldVal.Set(reflect.ValueOf(boolValue))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue := value.Int()
			if intValue == 0 {
				fieldVal.Set(reflect.ValueOf(false))
			} else {
				fieldVal.Set(reflect.ValueOf(true))
			}
		case reflect.Float32, reflect.Float64:
			floatValue := value.Float()
			if floatValue == 0 {
				fieldVal.Set(reflect.ValueOf(false))
			} else {
				fieldVal.Set(reflect.ValueOf(true))
			}
		default:
			return fmt.Errorf("unsupported type conversion: %s to %s", value.Kind(), fieldKind)
		}
	default:
		return fmt.Errorf("unsupported type conversion: %s to %s", value.Kind(), fieldKind)
	}

	return nil
}
