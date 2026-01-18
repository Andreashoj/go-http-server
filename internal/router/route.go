package router

import (
	"github.com/Andreashoj/go-http-server/internal/parser"
	"github.com/Andreashoj/go-http-server/internal/tests"
)

type route struct {
	Url     string
	Method  tests.Request // move away from tests package
	Handler func(writer HTTPWriter)
	Request *parser.HTTPRequest
}

func (r *router) Post(url string, handler func(writer HTTPWriter)) {
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  tests.Post,
	}
	r.add(newRoute)
}
