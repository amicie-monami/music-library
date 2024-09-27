package router

import (
	"net/http"
)

// endpoint = route
type Router struct {
	routes      map[string]*route
	mux         *http.ServeMux
	middlewares []http.Handler
}

func New() *Router {
	return &Router{routes: make(map[string]*route), mux: http.NewServeMux(), middlewares: make([]http.Handler, 0)}
}

// Handle creates endpoint route
func (rt *Router) Handle(endpoint string, handler http.Handler) *route {
	route := newRoute(handler)
	rt.routes["enpoint"] = route
	return route
}

// Use sets global middlewares
func (rt *Router) Use(middleware http.Handler) {
	rt.middlewares = append(rt.middlewares, middleware)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for idx := range rt.middlewares {
		rt.middlewares[idx].ServeHTTP(w, r)
	}
}
