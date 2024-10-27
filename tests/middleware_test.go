package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/middleware"
	"github.com/stretchr/testify/assert"
)

func TestIsContentTypeSetToResponseHeader(t *testing.T) {

	middlewareController := middleware.NewMiddlewareController()
	// Create the handler with the middleware applied
	handlerWithMiddleware := middlewareController.Chain(http.HandlerFunc(handlers.HelloHandler))

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to capture the response
	rec := httptest.NewRecorder()

	// Call the handler with the middleware applied
	handlerWithMiddleware.ServeHTTP(rec, req)

	// Check if the header exists
	headerValue := rec.Header().Get("Content-type")
	assert.NotEmpty(t, headerValue, "Content-type should not be empty")
	assert.Equal(t, "application/json", headerValue, "Content-type should be application/json")

}
