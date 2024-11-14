package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.HandlerFunc

type MiddleWareController struct {
	stack []Middleware
}

func NewMiddlewareController() *MiddleWareController {
	stack := []Middleware{
		SetHeaders,
		LogRequest,
	}

	return &MiddleWareController{
		stack: stack,
	}
}

func (mc MiddleWareController) Chain() Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(mc.stack) - 1; i >= 0; i-- {
			next = mc.stack[i](next)
		}

		return next.ServeHTTP
	}
}
