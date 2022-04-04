package core

import (
	"net/http"
)

type router struct {
	roots    map[string]*trie
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    map[string]*trie{},
		handlers: map[string]HandlerFunc{},
	}
}

// addRoute with request method and path to trie
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trie{
			son: map[string]*trie{},
		}
	}
	root := r.roots[method]
	root.insert(path)
	key := method + "-" + path
	r.handlers[key] = handler
}

// getRoute from the trie with request method
func (r *router) getRoute(method, path string) (*trie, map[string]string) {
	if root, ok := r.roots[method]; !ok {
		return nil, nil
	} else {
		return root.search(path)
	}
}

func (r *router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = params
		key := ctx.Method + "-" + node.path
		if handler, ok := r.handlers[key]; ok {
			handler(ctx)
		}
	} else {
		ctx.String(http.StatusNotFound, "404 not found")
	}
}
