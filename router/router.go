package router

import (
	"context"
	"net/http"
	"strings"
)

type paramsKey struct{}

type Router struct {
	root *node
}

func New() *Router {
	return &Router{
		&node{},
	}
}

func Param(r *http.Request, name string) string {
	params, _ := r.Context().Value(paramsKey{}).(map[string]string)
	return params[name]
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	current := r.root
	params := make(map[string]string)

	for _, p := range paths {
		if child, ok := current.staticChilder[p]; ok {
			current = child
		} else if current.paramChild != nil {
			params[current.paramChild.name] = p
			current = current.paramChild
		} else {
			// no match, so 404
			http.NotFound(w, req)
			return
		}
	}

	// current is final now, get handler
	handler, ok := current.handlers[req.Method]
	if !ok {
		http.NotFound(w, req)
		return
	}

	// inject params into context and call handler
	ctx := context.WithValue(req.Context(), paramsKey{}, params)
	handler(w, req.WithContext(ctx))
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	paths := strings.Split(strings.Trim(path, "/"), "/")

	current := r.root

	for _, p := range paths {
		if strings.HasPrefix(p, ":") {
			if current.paramChild == nil {
				current.paramChild = &node{name: p[1:]} // take out the :
			}
			current = current.paramChild // move node
		} else {
			if current.staticChilder == nil {
				current.staticChilder = make(map[string]*node)
			}
			_, ok := current.staticChilder[p]
			if !ok {
				current.staticChilder[p] = &node{}
			}
			current = current.staticChilder[p] // move
		}
	}

	// final node for handler
	if current.handlers == nil {
		current.handlers = make(map[string]http.HandlerFunc)
	}
	current.handlers[method] = handler

}
