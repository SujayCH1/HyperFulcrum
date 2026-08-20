package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type requestPayload struct {
	Name string `json:"name"`
}

func TestReadJSONRequest(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		expected    error
	}{
		{
			name:        "valid request",
			contentType: "application/json; charset=utf-8",
			body:        `{"name":"project"}`,
		},
		{
			name:        "unsupported content type",
			contentType: "text/plain",
			body:        `{"name":"project"}`,
			expected:    ErrUnsupportedMediaType,
		},
		{
			name:        "unknown field",
			contentType: "application/json",
			body:        `{"name":"project","unknown":true}`,
			expected:    errors.New("json: unknown field \"unknown\""),
		},
		{
			name:        "multiple objects",
			contentType: "application/json",
			body:        `{"name":"one"}{"name":"two"}`,
			expected:    errors.New("only one JSON object is allowed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()

			var payload requestPayload
			err := ReadJSONRequest(response, request, &payload)

			if test.expected == nil && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if test.expected != nil && (err == nil || err.Error() != test.expected.Error()) {
				t.Fatalf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func TestReadJSONRequestRejectsLargeBody(t *testing.T) {
	body := `{"name":"` + strings.Repeat("a", maxRequestBodySize) + `"}`
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	var payload requestPayload
	err := ReadJSONRequest(response, request, &payload)
	if !errors.Is(err, ErrRequestBodyTooLarge) {
		t.Fatalf("expected request body too large error, got %v", err)
	}
}

func TestWriteJSONErrorResponseRedactsInternalError(t *testing.T) {
	response := httptest.NewRecorder()

	WriteJSONErrorResponse(
		response,
		http.StatusInternalServerError,
		"Request failed",
		errors.New("database password leaked"),
	)

	if strings.Contains(response.Body.String(), "database password leaked") {
		t.Fatal("internal error was included in the response")
	}
}
