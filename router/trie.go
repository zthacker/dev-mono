package router

import "net/http"

type node struct {
	staticChilder map[string]*node
	paramChild    *node
	handlers      map[string]http.HandlerFunc
	name          string
}
