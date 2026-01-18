package router

import (
	"github.com/Andreashoj/go-http-server/internal/parser"
	"github.com/Andreashoj/go-http-server/internal/serializer"
	"github.com/Andreashoj/go-http-server/internal/tests"
)

type Router interface {
	Post(url string, handler func(writer serializer.HTTPWriter))
	FindMatchingRoute(request *parser.HTTPRequest) route
	add(route)
}

type router struct {
	routes []route
}

type route struct {
	Method  tests.Request // move away from tests package
	Handler func(writer serializer.HTTPWriter)
	Request *parser.HTTPRequest
}

func NewRouter() Router {
	return &router{}
}

func (r *router) add(route route) {
	r.routes = append(r.routes, route)
}

func (r *router) Post(url string, handler func(writer serializer.HTTPWriter)) {
	newRoute := route{
		Handler: handler,
		Method:  tests.Post,
	}
	r.add(newRoute)
}

func (r *router) FindMatchingRoute(request *parser.HTTPRequest) route {
	return r.routes[0]
}
