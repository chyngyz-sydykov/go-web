package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chyngyz-sydykov/go-web/application/middleware"
)

func TestIsContentTypeSetToResponseHeader(t *testing.T) {
	// arrange
	middlewareController := middleware.NewMiddlewareController()

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, Middleware!"))
	})

	handler := middlewareController.Chain()(finalHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// act
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// assert
	if resp.Header.Get("Content-type") != "application/json" {
		t.Errorf("Content-type header is not set, got '%s'", resp.Header.Get("Content-type"))
	}
}
