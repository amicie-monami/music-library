package middleware

import (
	"log/slog"
	"net/http"
)

// Log middleware logs the request metadata
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("request", "method", r.Method, "path", r.URL.Path, "protocol", r.Proto)
		next.ServeHTTP(w, r)
	})
}
