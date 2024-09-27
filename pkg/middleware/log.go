package middleware

import "net/http"

func Log() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return next
	}
}
