package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNodeConnectionValidatorStoresPayload(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`{"host":"localhost","port":5432,"database_name":"postgres","username":"postgres","password":"secret"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	called := false
	handler := NodeConnectionValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := GetNodeConnectionPayload(r)
		if !ok {
			t.Fatal("node connection payload was not stored")
		}
		if payload.Password != "secret" {
			t.Fatal("unexpected node connection payload")
		}
		called = true
	}))

	handler.ServeHTTP(response, request)

	if !called {
		t.Fatal("next handler was not called")
	}
}

func TestUUIDPathValidator(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		status int
	}{
		{name: "valid UUID", id: uuid.New().String(), status: http.StatusNoContent},
		{name: "invalid UUID", id: "invalid", status: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.SetPathValue("id", test.id)
			response := httptest.NewRecorder()

			handler := UUIDPathValidator("id")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			handler.ServeHTTP(response, request)

			if response.Code != test.status {
				t.Fatalf("expected status %d, got %d", test.status, response.Code)
			}
		})
	}
}

func TestRequestIDAndRecovery(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()
	handler := RequestID(Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})))

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}

	requestID := response.Header().Get("X-Request-ID")
	if _, err := uuid.Parse(requestID); err != nil {
		t.Fatalf("expected valid request ID, got %q", requestID)
	}

	if !strings.Contains(response.Body.String(), requestID) {
		t.Fatal("error response does not include request ID")
	}
}
