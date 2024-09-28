package router

import (
	"net/http"
)

type route struct {
	handler http.Handler
	methods map[string]struct{}
}

func newRoute(handler http.Handler) *route {
	return &route{methods: make(map[string]struct{}, 0), handler: handler}
}

func (rt *route) Methods(methods ...string) {
	for idx := range methods {
		method := methods[idx]
		rt.methods[method] = struct{}{}
	}
}

func (rt *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !rt.checkMethod(r.Method) {
		w.Write([]byte("Method not allowed"))
		return
	}
	rt.handler.ServeHTTP(w, r)
}

func (rt *route) checkMethod(method string) bool {
	_, ok := rt.methods[method]
	return ok
}
