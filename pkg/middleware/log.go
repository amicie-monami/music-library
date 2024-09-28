package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Log middleware logs the request metadata
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("new request, %s %s [%s]", r.Method, r.URL.Path, r.Proto))
		next.ServeHTTP(w, r)
	})
}
