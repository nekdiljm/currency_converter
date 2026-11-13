package middleware

import (
	"currency-converter/internal/api/response"
	"fmt"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				err := fmt.Errorf("panic: %v\n%s", p, debug.Stack())
				response.Error(w, err)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
