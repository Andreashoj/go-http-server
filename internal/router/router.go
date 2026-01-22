package router

import (
	"fmt"
	"strings"
)

type Router interface {
	Get(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Post(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Put(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Delete(url string, handler func(writer HTTPWriter, request HTTPRequest))
	FindMatchingRoute(request HTTPRequest) *route
	//Group(url string, router func(router Router))
	add(route)
}

type route struct {
	Url     string
	Method  Request // move away from tests package
	Handler func(writer HTTPWriter, request HTTPRequest)
	Request HTTPRequest
}

type router struct {
	routes []route // make this into a map that is easier to use for mapping
}

func NewRouter() Router {
	return &router{}
}

func (r *router) add(route route) {
	if len(route.Url) == 0 {
		panic(fmt.Sprintf("route must not be an empty string"))
	}

	if string(route.Url[len(route.Url)-1]) == "/" {
		panic(fmt.Sprintf("failed adding route, due to trailing slash"))
	}

	if string(route.Url[0]) != "/" {
		panic(fmt.Sprintf("failed adding route, should start with a /"))
	}

	r.routes = append(r.routes, route)
}

func (r *router) FindMatchingRoute(request HTTPRequest) *route {
	for _, routeEntry := range r.routes {
		if routeEntry.Method != request.Method() {
			continue
		}

		requestUrl := strings.Split(strings.TrimPrefix(request.Url(), "/"), "/") // Split route and request url and also removing leading / to avoid getting an empty entry in the slice
		routeUrl := strings.Split(strings.TrimPrefix(routeEntry.Url, "/"), "/")
		matches := true
		if len(requestUrl) != len(routeUrl) { // no need to check if route matches if their length isn't the same
			continue
		}

		for i, entry := range routeUrl {
			if string(entry[0]) == ":" { // dynamic route check ":id", any existent value from the request is accepted
				continue
			}

			if requestUrl[i] != entry { // not a dynamic route and the route param doesn't match
				matches = false
				break
			}
		}

		if matches {
			return &routeEntry
		}
	}

	return nil
}

func (r *router) Delete(url string, handler func(writer HTTPWriter, request HTTPRequest)) {
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Delete,
	}

	r.add(newRoute)
}

func (r *router) Put(url string, handler func(writer HTTPWriter, request HTTPRequest)) {
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Put,
	}

	r.add(newRoute)
}

func (r *router) Post(url string, handler func(writer HTTPWriter, request HTTPRequest)) {
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Post,
	}

	r.add(newRoute)
}

func (r *router) Get(url string, handler func(writer HTTPWriter, request HTTPRequest)) {
	fmt.Println("url", url)
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Get,
	}

	r.add(newRoute)
}
