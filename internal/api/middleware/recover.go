package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				err := fmt.Errorf("panic: %v\n%s", p, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
