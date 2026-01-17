package router

import (
	"net"

	"github.com/Andreashoj/go-http-server/internal/parser"
)

type Router interface {
	Post(url string, handler func(cn net.Conn))
	FindMatchingRoute(request *parser.HTTPRequest) route
	add(route)
}

type router struct {
	routes []route
}

type route struct {
	Handler func(cn net.Conn)
}

func NewRouter() Router {
	return &router{}
}

func (r *router) add(route route) {
	r.routes = append(r.routes, route)
}

func (r *router) Post(url string, handler func(cn net.Conn)) {
	newRoute := route{
		Handler: handler,
	}
	r.add(newRoute)
}

func (r *router) FindMatchingRoute(request *parser.HTTPRequest) route {
	return r.routes[0]
}
