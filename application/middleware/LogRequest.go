package middleware

import (
	"fmt"
	"net/http"
)

func LogRequest(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("method %s\n", r.Method)
		next.ServeHTTP(w, r)
	})
}
