package rag

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
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
