package middleware

import (
	"net/http"
	"strings"
)

func EnforceAccept(requiredType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept")

			if !strings.Contains(accept, requiredType) && !strings.Contains(accept, "*/*") {
				http.Error(w, "Not Acceptable: Requires "+requiredType, http.StatusNotAcceptable)
				return
			}

			next.ServeHTTP(w,r)
		})
	}
}
