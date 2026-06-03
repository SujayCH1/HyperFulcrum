package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(payload)
}

func WriteJSONSuccessResponse(w http.ResponseWriter,status int,message string,data any) error {

	response := SuccessResponse{Success: true, Message: message, Data: data}

	return WriteJSONResponse(w,status,response)
}

func WriteJSONErrorResponse(w http.ResponseWriter,status int,message string,err error) error {

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	response := ErrorResponse{Success: false, Message: message, Error: errStr}

	return WriteJSONResponse(w,status,response)
}

func ReadJSONRequest(r *http.Request, result any) error {

	if r.Body == nil {
		return errors.New("request body is required")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields() //rejects unknown fields in the JSON payolad

	if err := decoder.Decode(result); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body cannot be empty")
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

	if decoder.More() {
		return errors.New(
			"only one JSON object is allowed",
		)
	}
	return nil
}
