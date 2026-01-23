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
	Group(url string, router func(router Router))
	add(route)
}

type route struct {
	Url     string
	Method  Request
	Handler func(writer HTTPWriter, request HTTPRequest)
	Request HTTPRequest
	Group   []router
}

type router struct {
	routes      []route
	currentNode *node
}

type node struct {
	parent   *node
	children []node
	path     string
	route    *route

	// configs
	middlewares any
	headers     any
}

func NewRouter() Router {
	nde := node{
		path: "/", // Initial route
	}
	return &router{
		currentNode: &nde,
	}
}

func (r *router) add(route route) {
	if len(route.Url) == 0 {
		panic(fmt.Sprintf("route must not be an empty string"))
	}

	if string(route.Url[len(route.Url)-1]) == "/" {
		panic(fmt.Sprintf("failed adding route, shouldn't end with a /"))
	}

	if string(route.Url[0]) != "/" {
		panic(fmt.Sprintf("failed adding route, should start with a /"))
	}

	nde := node{
		parent: r.currentNode,
		path:   route.Url,
		route:  &route,
	}

	r.currentNode.children = append(r.currentNode.children, nde)
	r.routes = append(r.routes, route) // This may not be needed
}

func findMatchingNode(requestUrl string, node *node) *node {
	isRoot := node.path == "/"
	if node.path == requestUrl {
		return node
	}

	var requestUrlsParts []string
	urlParts := strings.Split(requestUrl, "/")
	for i, part := range urlParts {
		if isRoot && i != 0 {
			requestUrlsParts = append(requestUrlsParts, "/"+part)
		} else if i != 0 && i != 1 { // removes current entry + leading /
			requestUrlsParts = append(requestUrlsParts, "/"+part)
		}
	}

	// Check again on current path
	for _, child := range node.children {
		if child.path == requestUrlsParts[0] {
			return findMatchingNode(strings.Join(requestUrlsParts, ""), &child)
		}
	}

	return nil
}

func (r *router) FindMatchingRoute(request HTTPRequest) *route {
	n := findMatchingNode(request.Url(), r.currentNode)
	return n.route
}

func (r *router) Group(url string, handler func(router Router)) {
	var groupUrl = "/" + url
	if r.currentNode.path == "/" {
		groupUrl = url
	}

	nde := node{
		parent:   r.currentNode,
		path:     groupUrl,
		children: []node{},
	}

	rter := router{
		currentNode: &nde,
	}

	fmt.Println()
	handler(&rter)
	r.currentNode.children = append(r.currentNode.children, nde)
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
