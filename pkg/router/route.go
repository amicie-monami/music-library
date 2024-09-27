package router

import "net/http"

// method = handler
type route struct {
	handler http.Handler
	methods []string
}

func newRoute(handler http.Handler) *route {
	return &route{methods: make([]string, 0), handler: handler}
}

func (r *route) Methods(methods ...string) {
	r.methods = append(r.methods, methods...)
}
