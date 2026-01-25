package router

func Respond(writer HTTPWriter, request HTTPRequest, middlewares []MiddlewareFunc, httpHandler func(writer HTTPWriter, request HTTPRequest)) {
	handler := func() {
		httpHandler(writer, request)
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewareWrapper(writer, request, handler, middlewares[i])
	}

	handler()
}

type MiddlewareFunc func(writer HTTPWriter, request HTTPRequest, next func())

func middlewareWrapper(w HTTPWriter, r HTTPRequest, next func(), middleware MiddlewareFunc) func() {
	return func() {
		middleware(w, r, next)
	}
}

func GetMiddlewares(node *node) []MiddlewareFunc {
	middlewares := node.middlewares

	if node.parent != nil {
		middlewares = append(middlewares, GetMiddlewares(node.parent)...)
	}

	return middlewares
}

func (r *router) Use(middlewareFunc func(writer HTTPWriter, request HTTPRequest, next func())) {
	r.currentNode.middlewares = append(r.currentNode.middlewares, middlewareFunc)
}
