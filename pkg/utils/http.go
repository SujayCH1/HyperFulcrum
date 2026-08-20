package utils

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxRequestBodySize = 1 << 20

var (
	ErrRequestBodyTooLarge  = errors.New("request body is too large")
	ErrUnsupportedMediaType = errors.New("content type must be application/json")
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(payload)
}

func WriteJSONSuccessResponse(w http.ResponseWriter, status int, message string, data any) error {

	response := SuccessResponse{Success: true, Message: message, Data: data}

	return WriteJSONResponse(w, status, response)
}

func WriteJSONErrorResponse(w http.ResponseWriter, status int, message string, err error) error {

	var errStr string
	if err != nil && status < http.StatusInternalServerError {
		errStr = err.Error()
	}

	response := ErrorResponse{
		Success:   false,
		Message:   message,
		Error:     errStr,
		RequestID: w.Header().Get("X-Request-ID"),
	}

	return WriteJSONResponse(w, status, response)
}

func ReadJSONRequest(w http.ResponseWriter, r *http.Request, result any) error {

	if r.Body == nil {
		return errors.New("request body is required")
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "application/json" {
		return ErrUnsupportedMediaType
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err = decoder.Decode(result); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body cannot be empty")
		}

		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return ErrRequestBodyTooLarge
		}

		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			return errors.New("invalid JSON syntax")
		}

		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			return errors.New(
				"invalid type for field: " +
					typeError.Field,
			)
		}
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return ErrRequestBodyTooLarge
		}

		return errors.New("only one JSON object is allowed")
	}

	return nil
}
