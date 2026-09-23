package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeDatabase struct{}

func (fakeDatabase) Ping(context.Context) error {
	return nil
}

type failingDatabase struct{}

func (failingDatabase) Ping(context.Context) error {
	return errors.New("database unavailable")
}

func TestCheckReturnsOK(t *testing.T) {
	handler := NewHandler(fakeDatabase{})

	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()

	handler.Check(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	expected := "{\"status\":\"ok\"}\n"

	if response.Body.String() != expected {
		t.Fatalf(
			"expected body %q, got %q",
			expected,
			response.Body.String(),
		)
	}
}

func TestCheckReturnsServiceUnavailableWhenDatabaseFails(t *testing.T) {
	handler := NewHandler(failingDatabase{})

	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()

	handler.Check(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			response.Code,
		)
	}

	expected := "{\"status\":\"unhealthy\"}\n"

	if response.Body.String() != expected {
		t.Fatalf(
			"expected body %q, got %q",
			expected,
			response.Body.String(),
		)
	}
}
