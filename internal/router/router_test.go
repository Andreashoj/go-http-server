package router

import (
	"testing"
)

func Test_router_FindMatchingRoute(t *testing.T) {
	routerMock := router{}
	routerMock.Get("/url", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Get("/url/:id", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Post("/users/example", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Put("/user", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Get("/posts/:postId/comments/:commentId", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Delete("/items/:id", func(writer HTTPWriter, request HTTPRequest) {})

	requests := []struct {
		endpoint    httpRequest
		expectedURL string
	}{
		{
			expectedURL: "/url/:id",
			endpoint: httpRequest{
				startLine: "GET /url/123 HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "GET",
			},
		},
		{
			expectedURL: "/url",
			endpoint: httpRequest{
				startLine: "GET /url HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/url",
				routerURL: "/url",
				method:    "GET",
			},
		},
		{
			expectedURL: "/users/example",
			endpoint: httpRequest{
				startLine: "POST /users/example HTTP/1.1",
				headers:   []string{"Host: example.com", "Content-Type: application/json"},
				body:      `{"name":"test"}`,
				params:    nil,
				url:       "/users/example",
				routerURL: "/users/example",
				method:    "POST",
			},
		},
		{
			expectedURL: "/user",
			endpoint: httpRequest{
				startLine: "PUT /user HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/user",
				routerURL: "/user",
				method:    "PUT",
			},
		},
		{
			expectedURL: "/posts/:postId/comments/:commentId",
			endpoint: httpRequest{
				startLine: "GET /posts/42/comments/789 HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/posts/42/comments/789",
				routerURL: "/posts/:postId/comments/:commentId",
				method:    "GET",
			},
		},
		{
			expectedURL: "/items/:id",
			endpoint: httpRequest{
				startLine: "DELETE /items/999 HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/items/999",
				routerURL: "/items/:id",
				method:    "DELETE",
			},
		},
		{
			expectedURL: "",
			endpoint: httpRequest{
				startLine: "GET /nonexistent HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/nonexistent",
				routerURL: "/nonexistent",
				method:    "GET",
			},
		},
		{
			expectedURL: "",
			endpoint: httpRequest{
				startLine: "POST /url/123 HTTP/1.1",
				headers:   []string{"Host: example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "POST",
			},
		},
	}

	for _, tt := range requests {
		routerEndpoint := routerMock.FindMatchingRoute(&tt.endpoint)
		if tt.expectedURL == "" {
			if routerEndpoint != nil {
				t.Errorf("failed, expected nil but got %s", routerEndpoint.Url)
			}
		} else {
			if routerEndpoint == nil {
				t.Errorf("test failed, expected %s but got nil", tt.expectedURL)
			}

			if routerEndpoint.Url != tt.expectedURL {
				t.Errorf("expected %s but got %s", tt.expectedURL, routerEndpoint.Url)
			}
		}

	}
}

func Test_router_AddRoutes(t *testing.T) {
	t.Run("valid routes", func(t *testing.T) {
		r := NewRouter()
		routes := []route{
			{
				Url:     "/users",
				Method:  Post,
				Handler: nil,
				Request: nil,
			},
			{
				Url:     "/users/:id",
				Method:  Put,
				Handler: nil,
				Request: nil,
			},
			{
				Url:     "/users/:id/example/test",
				Method:  Get,
				Handler: nil,
				Request: nil,
			},
		}

		for _, tt := range routes {
			assertNoPanic(t, func() {
				r.add(tt)
			})
		}
	})

	t.Run("invalid routes", func(t *testing.T) {
		r := NewRouter()
		routes := []route{
			{
				Url:     "/users/",
				Method:  Post,
				Handler: nil,
				Request: nil,
			},
			{
				Url:     "",
				Method:  Put,
				Handler: nil,
				Request: nil,
			},
			{
				Url:     "users/:id/example/test",
				Method:  Get,
				Handler: nil,
				Request: nil,
			},
		}

		for _, tt := range routes {
			assertPanic(t, func() {
				r.add(tt)
			})
		}
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func assertNoPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	f()
}
