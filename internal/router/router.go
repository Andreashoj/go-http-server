package router

import (
	"github.com/Andreashoj/go-http-server/internal/parser"
)

type Router interface {
	Post(url string, handler func(writer HTTPWriter))
	FindMatchingRoute(request *parser.HTTPRequest) route
	add(route)
}

type router struct {
	routes []route
}

func NewRouter() Router {
	return &router{}
}

func (r *router) add(route route) {
	r.routes = append(r.routes, route)
}

func (r *router) FindMatchingRoute(request *parser.HTTPRequest) route {
	return r.routes[0]
}
