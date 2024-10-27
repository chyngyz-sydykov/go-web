package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MiddleWareController struct {
	stack []Middleware
}

func NewMiddlewareController() *MiddleWareController {
	stack := []Middleware{
		SetHeadersMiddleware,
	}

	return &MiddleWareController{
		stack: stack,
	}
}

func (mc MiddleWareController) Chain(h http.HandlerFunc) http.HandlerFunc {
	if len(mc.stack) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(mc.stack) - 1; i >= 0; i-- {
		wrapped = mc.stack[i](wrapped)
	}

	return wrapped
}
