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
	FindMatchingRoute(request HTTPRequest) (*route, error)
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
}

func compareRoutes(requestUrl, routerUrl string) bool {
	requestUrlParts := strings.Split(requestUrl, "/")
	routerUrlParts := strings.Split(routerUrl, "/")

	if routerUrl == "/" {
		routerUrlParts = []string{""}
	}

	if requestUrl == "/" {
		requestUrlParts = []string{""}
	}

	if len(requestUrlParts) != len(routerUrlParts) {
		return false
	}

	for i, part := range routerUrlParts {
		if part == "" {
			continue
		}

		isDynamic := string(part[0]) == ":"
		if isDynamic {
			continue
		}

		if part != requestUrlParts[i] {
			return false
		}
	}

	return true
}

func findMatchingNode(requestUrl string, method Request, n *node) *node {
	isRoot := n.path == "/"
	if compareRoutes(requestUrl, n.path) && n.route.Method == method {
		return n
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
	var result *node
	for _, child := range n.children {
		if compareRoutes(requestUrl, child.path) && child.route.Method == method {
			result = &child
		} else if child.path == requestUrlsParts[0] {
			result = findMatchingNode(strings.Join(requestUrlsParts, ""), method, &child)
		}

		if result != nil {
			return result
		}
	}

	return nil
}

func (r *router) FindMatchingRoute(request HTTPRequest) (*route, error) {
	n := findMatchingNode(request.Url(), request.Method(), r.currentNode)
	if n == nil {
		return nil, fmt.Errorf("could not find match for request URL: %s", request.Url())
	}
	return n.route, nil
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
