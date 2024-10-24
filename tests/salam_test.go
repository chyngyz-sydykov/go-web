package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	assert := assert.New(t)

	assert.HTTPStatusCode(HelloHandler, "GET", "/", nil, 200)
	// if err != nil {
	// 	t.Fatalf("Could not create request: %v", err)
	// }

	// // Check the status code is what we expect
	// if rec.Code != http.StatusOK {
	// 	t.Errorf("Expected status OK; got %v", rec.Code)
	// }

	// // Check the response body is what we expect
	// expected := "hello world"
	// if rec.Body.String() != expected {
	// 	t.Errorf("Expected body to be %q; got %q", expected, rec.Body.String())
	// }
}
