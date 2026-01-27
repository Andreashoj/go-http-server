# go-http-server

A lightweight HTTP router library for Go. Build simple HTTP servers without the bloat.

## Install

```bash
go get github.com/Andreashoj/go-http-server
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/Andreashoj/go-http-server/router"
)

func main() {
	r := router.New()

	// Define a route
	r.Get("/hello", func(w router.HTTPWriter, req router.HTTPRequest) {
		w.Header().Add(router.ContentType, "text/plain")
		w.FormatResponse("Hello, World!", 200)
	})

	// Start listening
	router.StartServer(":8080", r)
}
```

Visit `http://localhost:8080/hello` and boom, you got a response.

## Features

- **Simple routing** — Define routes with HTTP methods (GET, POST, PUT, DELETE, etc.)
- **Middleware support** — Chain middleware through your routes
- **Route nesting** — Organize routes hierarchically
- **URL parameters** — Extract dynamic params from routes

## Examples

### POST with JSON
```go
r.Post("/api/users", func(w router.HTTPWriter, req router.HTTPRequest) {
	// Get the request body
	body := req.Body()
	
	w.Header().Add(router.ContentType, "application/json")
	w.FormatResponse(`{"id":1,"name":"test"}`, 201)
})
```

### Extract URL Parameters
```go
r.Get("/users/:id", func(w router.HTTPWriter, req router.HTTPRequest) {
	id, err := req.GetURLParam("id")
	if err != nil {
		w.FormatResponse("User not found", 404)
		return
	}
	
	w.Header().Add(router.ContentType, "application/json")
	w.FormatResponse(fmt.Sprintf(`{"id":"%s"}`, id), 200)
})
```

### Query Parameters
```go
r.Get("/search", func(w router.HTTPWriter, req router.HTTPRequest) {
	query, err := req.GetQueryParam("q")
	if err != nil {
		w.FormatResponse("Missing query parameter", 400)
		return
	}
	
	w.Header().Add(router.ContentType, "application/json")
	w.FormatResponse(fmt.Sprintf(`{"results":"%s"}`, query), 200)
})
```

### Middleware
```go
authMiddleware := func(w router.HTTPWriter, req router.HTTPRequest, next func()) {
    // Check auth headers or tokens
    if req.Method() != router.Get {
        next()
        return
    }
	
    next()
}

r.Use(authMiddleware)
r.Get("/protected", handler)
```

### Nested Routes
```go
api := r.Group("/api")
api.Get("/users", getUsers)
api.Post("/users", createUser)
```

## License

MIT