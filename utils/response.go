package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/devjoemedia/scrumpilot-go-api/types"
)

func JSON[T any](w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Response encoding failed", http.StatusInternalServerError)
	}
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// 1. Limit request body size to 1MB to prevent DOS attacks
	maxBytes := int64(1048576)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)

	// 2. Optional: Return an error if JSON contains fields not in your struct
	decoder.DisallowUnknownFields()

	// 3. Decode the body into the destination
	if err := decoder.Decode(data); err != nil {
		return err
	}

	// 4. Ensure there is only ONE JSON object in the request body
	err := decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := types.ErrorResponse{
		Message: msg,
		Success: false,
		Status:  status,
	}
	// json.NewEncoder(w).Encode(resp)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
	}
}
