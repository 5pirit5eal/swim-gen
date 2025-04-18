package models

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"reflect"
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

// StructToMap converts a struct to a map[string]any.
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
