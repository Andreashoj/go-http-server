package router

import (
	"testing"
)

func Test_router_FindMatchingRoute(t *testing.T) {
	routerMock := NewRouter()
	routerMock.Get("/url", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Get("/url/:id", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Post("/users/example", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Put("/user", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Get("/posts/:postId/comments/:commentId", func(writer HTTPWriter, request HTTPRequest) {})
	routerMock.Delete("/items/:id", func(writer HTTPWriter, request HTTPRequest) {})

	tests := []struct {
		name        string
		request     httpRequest
		expectedURL string
	}{
		{
			name: "exact match static route",
			request: httpRequest{
				startLine: "GET /url HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/url",
				routerURL: "/url",
				method:    "GET",
			},
			expectedURL: "/url",
		},
		{
			name: "dynamic parameter match",
			request: httpRequest{
				startLine: "GET /url/123 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "GET",
			},
			expectedURL: "/url/:id",
		},
		{
			name: "static POST route",
			request: httpRequest{
				startLine: "POST /users/example HTTP/1.1",
				headers:   map[string]string{"Host": "example.com", "Content-Type": "application/json"},
				body:      `{"name":"test"}`,
				params:    nil,
				url:       "/users/example",
				routerURL: "/users/example",
				method:    "POST",
			},
			expectedURL: "/users/example",
		},
		{
			name: "static PUT route",
			request: httpRequest{
				startLine: "PUT /user HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/user",
				routerURL: "/user",
				method:    "PUT",
			},
			expectedURL: "/user",
		},
		{
			name: "multiple dynamic parameters",
			request: httpRequest{
				startLine: "GET /posts/42/comments/789 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/posts/42/comments/789",
				routerURL: "/posts/:postId/comments/:commentId",
				method:    "GET",
			},
			expectedURL: "/posts/:postId/comments/:commentId",
		},
		{
			name: "DELETE route with dynamic parameter",
			request: httpRequest{
				startLine: "DELETE /items/999 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/items/999",
				routerURL: "/items/:id",
				method:    "DELETE",
			},
			expectedURL: "/items/:id",
		},
		{
			name: "nonexistent route returns nil",
			request: httpRequest{
				startLine: "GET /nonexistent HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/nonexistent",
				routerURL: "/nonexistent",
				method:    "GET",
			},
			expectedURL: "",
		},
		{
			name: "wrong method returns nil",
			request: httpRequest{
				startLine: "POST /url/123 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "POST",
			},
			expectedURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routerNode, err := routerMock.FindMatchingRoute(&tt.request)

			if tt.expectedURL == "" {
				if err == nil {
					t.Errorf("expected error but got %s", routerNode.Route.Url)
				}
			} else {
				if routerNode == nil {
					t.Errorf("expected %s but got nil", tt.expectedURL)
				} else if routerNode.Route.Url != tt.expectedURL {
					t.Errorf("expected %s but got %s", tt.expectedURL, routerNode.Route.Url)
				}
			}
		})
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

func TestCompareRoutes(t *testing.T) {
	tests := []struct {
		name       string
		requestUrl string
		routerUrl  string
		expected   bool
	}{
		// Exact matches
		{
			name:       "exact match",
			requestUrl: "/users/profile",
			routerUrl:  "/users/profile",
			expected:   true,
		},
		{
			name:       "exact match single segment",
			requestUrl: "/users",
			routerUrl:  "/users",
			expected:   true,
		},
		// Dynamic parameters
		{
			name:       "dynamic parameter match",
			requestUrl: "/users/123",
			routerUrl:  "/users/:id",
			expected:   true,
		},
		{
			name:       "multiple dynamic parameters",
			requestUrl: "/users/123/posts/456",
			routerUrl:  "/users/:userId/posts/:postId",
			expected:   true,
		},
		{
			name:       "mixed static and dynamic",
			requestUrl: "/api/users/john/profile",
			routerUrl:  "/api/users/:name/profile",
			expected:   true,
		},
		// Length mismatches
		{
			name:       "different number of segments",
			requestUrl: "/users/123",
			routerUrl:  "/users/123/posts",
			expected:   false,
		},
		{
			name:       "request has fewer segments",
			requestUrl: "/users",
			routerUrl:  "/users/123/posts",
			expected:   false,
		},
		{
			name:       "router has fewer segments",
			requestUrl: "/users/123/posts",
			routerUrl:  "/users/123",
			expected:   false,
		},
		// Static mismatch
		{
			name:       "static segments don't match",
			requestUrl: "/users/123",
			routerUrl:  "/posts/123",
			expected:   false,
		},
		{
			name:       "different static segment in middle",
			requestUrl: "/api/users/123/profile",
			routerUrl:  "/api/posts/123/profile",
			expected:   false,
		},
		// Edge cases
		{
			name:       "all dynamic parameters",
			requestUrl: "/a/b/c",
			routerUrl:  "/:x/:y/:z",
			expected:   true,
		},
		{
			name:       "all static segments",
			requestUrl: "/users/profile/settings",
			routerUrl:  "/users/profile/settings",
			expected:   true,
		},
		{
			name:       "empty paths",
			requestUrl: "/",
			routerUrl:  "/",
			expected:   true,
		},
		{
			name:       "single character dynamic parameter",
			requestUrl: "/a",
			routerUrl:  "/:id",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareRoutes(tt.requestUrl, tt.routerUrl)
			if result != tt.expected {
				t.Errorf("compareRoutes(%q, %q) = %v, want %v",
					tt.requestUrl, tt.routerUrl, result, tt.expected)
			}
		})
	}
}
