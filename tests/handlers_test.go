package main

import (
	"testing"

	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	assert := assert.New(t)

	assert.HTTPStatusCode(handlers.HelloHandler, "GET", "/", nil, 200)
	assert.HTTPBodyContains(handlers.HelloHandler, "GET", "/", nil, "hello world")

}
