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
