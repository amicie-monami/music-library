package router

import (
	"net/http"
)

type Router struct {
	routes map[string]*route
	mux    *http.ServeMux
}

func New() *Router {
	return &Router{routes: make(map[string]*route), mux: http.NewServeMux()}
}

// Handle creates endpoint route
func (rt *Router) Handle(endpoint string, handler http.Handler) *route {
	rt.routes[endpoint] = newRoute(handler)
	return rt.routes[endpoint]
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := rt.routes[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	route.ServeHTTP(w, r)
}
