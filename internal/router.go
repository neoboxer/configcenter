package internal

import (
	"net/http"
)

type Router struct {
	server *internalServer
}

func NewRouter() *Router {
	server := &internalServer{}
	return &Router{server: server}
}

// func (r *Router) register(method string, pattern string, handleFunc) {
// }

func (r *Router) Run(addr string) error {
	http.ListenAndServe(addr, r.server)
	return nil
}

// func (r *Router) GET(pattern string, handlerFunc http.HandlerFunc) {
// }
