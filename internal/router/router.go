package router

type Router interface {
	Get(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Post(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Put(url string, handler func(writer HTTPWriter, request HTTPRequest))
	Delete(url string, handler func(writer HTTPWriter, request HTTPRequest))
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
	routes []route // make this into a map that is easier to use for mapping
}

func NewRouter() Router {
	return &router{}
}

func (r *router) add(route route) {
	r.routes = append(r.routes, route)
}

func (r *router) FindMatchingRoute(request HTTPRequest) *route {
	for _, routeEntry := range r.routes {
		if routeEntry.Method == request.Method() && routeEntry.Url == request.Url() {
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
	newRoute := route{
		Url:     url,
		Handler: handler,
		Method:  Get,
	}

	r.add(newRoute)
}
