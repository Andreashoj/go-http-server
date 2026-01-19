package router

type Router interface {
	Post(url string, handler func(writer HTTPWriter, request HTTPRequest))
	FindMatchingRoute(request HTTPRequest) *route
	add(route)
}

type route struct {
	Url     string
	Method  Request // move away from tests package
	Handler func(writer HTTPWriter, request HTTPRequest)
	Request HTTPRequest
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

func (r *router) FindMatchingRoute(request HTTPRequest) *route {
	for _, routeEntry := range r.routes {
		if routeEntry.Url == request.Url() { // Make sure that
			return &routeEntry
		}
	}

	return nil
}

func (r *router) Post(url string, handler func(writer HTTPWriter, request HTTPRequest)) {
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Post,
	}

	r.add(newRoute)
}
