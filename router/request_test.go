package router

import "testing"

// Test get url param
func Test_httpRequest_GetQueryParam(t *testing.T) {
	requests := []struct {
		request       httpRequest
		key           string
		expectedValue string
	}{
		{
			key:           "test",
			expectedValue: "123",
			request: httpRequest{
				startLine: "GET /url?test=123 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    map[string]string{"test": "123"},
				url:       "/url",
				routerURL: "/url",
				method:    "GET",
			},
		},
		{
			key:           "id",
			expectedValue: "456",
			request: httpRequest{
				startLine: "GET /users?id=456&name=john HTTP/1.1",
				headers:   map[string]string{"Host": "api.example.com"},
				body:      "",
				params:    map[string]string{"id": "456", "name": "john"},
				url:       "/users",
				routerURL: "/users",
				method:    "GET",
			},
		},
		{
			key:           "search",
			expectedValue: "golang",
			request: httpRequest{
				startLine: "POST /search HTTP/1.1",
				headers:   map[string]string{"Host": "example.com", "Content-Type": "application/x-www-form-urlencoded"},
				body:      "search=golang&limit=10",
				params:    map[string]string{"search": "golang", "limit": "10"},
				url:       "/search",
				routerURL: "/search",
				method:    "POST",
			},
		},
		{
			key:           "nonexistent",
			expectedValue: "",
			request: httpRequest{
				startLine: "GET /test HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    map[string]string{},
				url:       "/test",
				routerURL: "/test",
				method:    "GET",
			},
		},
		{
			key:           "token",
			expectedValue: "abc123xyz",
			request: httpRequest{
				startLine: "GET /api/data?token=abc123xyz HTTP/1.1",
				headers:   map[string]string{"Host": "api.example.com", "Authorization": "Bearer token"},
				body:      "",
				params:    map[string]string{"token": "abc123xyz"},
				url:       "/api/data",
				routerURL: "/api/data",
				method:    "GET",
			},
		},
	}

	for _, tt := range requests {
		val, err := tt.request.GetQueryParam(tt.key)
		if tt.expectedValue != "" && err != nil {
			t.Errorf("test failed: %s", err)
		}

		if tt.expectedValue != val {
			t.Errorf("expected value %s but got %s", tt.expectedValue, val)
		}
	}
}

func Test_httpRequest_GetURLParam(t *testing.T) {
	requests := []struct {
		request       httpRequest
		key           string
		expectedValue string
	}{
		{
			key:           "id",
			expectedValue: "123",
			request: httpRequest{
				startLine: "GET /url/123 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "GET",
			},
		},
		{
			key:           "username",
			expectedValue: "john",
			request: httpRequest{
				startLine: "GET /users/john HTTP/1.1",
				headers:   map[string]string{"Host": "api.example.com"},
				body:      "",
				params:    nil,
				url:       "/users/john",
				routerURL: "/users/:username",
				method:    "GET",
			},
		},
		{
			key:           "id",
			expectedValue: "456",
			request: httpRequest{
				startLine: "GET /posts/456/comments/789 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/posts/456/comments/789",
				routerURL: "/posts/:id/comments/:commentId",
				method:    "GET",
			},
		},
		{
			key:           "commentId",
			expectedValue: "789",
			request: httpRequest{
				startLine: "GET /posts/456/comments/789 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/posts/456/comments/789",
				routerURL: "/posts/:id/comments/:commentId",
				method:    "GET",
			},
		},
		{
			key:           "nonexistent",
			expectedValue: "",
			request: httpRequest{
				startLine: "GET /url/123 HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/url/123",
				routerURL: "/url/:id",
				method:    "GET",
			},
		},
		{
			key:           "slug",
			expectedValue: "my-post-title",
			request: httpRequest{
				startLine: "GET /blog/my-post-title HTTP/1.1",
				headers:   map[string]string{"Host": "example.com"},
				body:      "",
				params:    nil,
				url:       "/blog/my-post-title",
				routerURL: "/blog/:slug",
				method:    "GET",
			},
		},
	}

	for _, tt := range requests {
		val, err := tt.request.GetURLParam(tt.key)
		if tt.expectedValue != "" && err != nil {
			t.Errorf("test failed: %s", err)
		}

		if tt.expectedValue != val {
			t.Errorf("expected value %s but got %s", tt.expectedValue, val)
		}
	}
}
