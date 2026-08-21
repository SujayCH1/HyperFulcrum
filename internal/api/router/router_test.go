package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hyperfulcrum/internal/api/handlers"
)

func newTestRouter() http.Handler {
	return NewRouter(
		handlers.NewProjectHandler(nil),
		handlers.NewNodeHandler(nil),
		handlers.NewNodeConnectionHandler(nil),
		handlers.NewTopoogyHandler(nil),
		handlers.NewShardKeyHandler(nil),
		handlers.NewReplicationHandler(nil),
	)
}

func TestReplicationRoutesReturnNotImplemented(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/replication", nil)
	response := httptest.NewRecorder()

	newTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d", response.Code)
	}
}

func TestRouterRejectsTrailingSlash(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/projects/", nil)
	response := httptest.NewRecorder()

	newTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestRouterReturnsPlainTextMethodNotAllowed(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/projects", nil)
	response := httptest.NewRecorder()

	newTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}

	if response.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatal("method not allowed response is not plain text")
	}

	if !strings.Contains(response.Header().Get("Allow"), http.MethodGet) {
		t.Fatal("method not allowed response does not include allowed methods")
	}
}

func TestRouterRejectsInvalidUUID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/projects/invalid", nil)
	response := httptest.NewRecorder()

	newTestRouter().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}
