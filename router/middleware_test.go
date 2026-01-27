package router

import (
	"reflect"
	"slices"
	"testing"
)

func TestApplyMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		middlewares []struct {
			pre  string
			post string
		}
		expectedExecutionOrder []string
		routeNode              *node
		middlewareContext      func(expExecutionOrder []string) []MiddlewareFunc
	}{
		{
			name: "single middleware",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "mw1-pre", post: "mw1-post"},
			},
			expectedExecutionOrder: []string{"mw1-pre", "mw1-post"},
		},
		{
			name: "two middlewares",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "mw1-pre", post: "mw1-post"},
				{pre: "mw2-pre", post: "mw2-post"},
			},
			expectedExecutionOrder: []string{"mw1-pre", "mw2-pre", "mw2-post", "mw1-post"},
		},
		{
			name: "three middlewares",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "mw1-pre", post: "mw1-post"},
				{pre: "mw2-pre", post: "mw2-post"},
				{pre: "mw3-pre", post: "mw3-post"},
			},
			expectedExecutionOrder: []string{"mw1-pre", "mw2-pre", "mw3-pre", "mw3-post", "mw2-post", "mw1-post"},
		},
		{
			name: "no middlewares",
			middlewares: []struct {
				pre  string
				post string
			}{},
			expectedExecutionOrder: []string{},
		},
		{
			name: "four middlewares in sequence",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "auth-pre", post: "auth-post"},
				{pre: "logging-pre", post: "logging-post"},
				{pre: "cors-pre", post: "cors-post"},
				{pre: "validate-pre", post: "validate-post"},
			},
			expectedExecutionOrder: []string{
				"auth-pre", "logging-pre", "cors-pre", "validate-pre",
				"validate-post", "cors-post", "logging-post", "auth-post",
			},
		},
		{
			name: "middlewares with nested route hierarchy",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "global-mw-pre", post: "global-mw-post"},
				{pre: "route-mw-pre", post: "route-mw-post"},
			},
			expectedExecutionOrder: []string{
				"global-mw-pre", "route-mw-pre", "route-mw-post", "global-mw-post",
			},
		},
		{
			name: "many middlewares",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "mw1-pre", post: "mw1-post"},
				{pre: "mw2-pre", post: "mw2-post"},
				{pre: "mw3-pre", post: "mw3-post"},
				{pre: "mw4-pre", post: "mw4-post"},
				{pre: "mw5-pre", post: "mw5-post"},
			},
			expectedExecutionOrder: []string{
				"mw1-pre", "mw2-pre", "mw3-pre", "mw4-pre", "mw5-pre",
				"mw5-post", "mw4-post", "mw3-post", "mw2-post", "mw1-post",
			},
		},
		{
			name: "authentication and logging middlewares",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "request-id-pre", post: "request-id-post"},
				{pre: "auth-check-pre", post: "auth-check-post"},
				{pre: "permission-check-pre", post: "permission-check-post"},
			},
			expectedExecutionOrder: []string{
				"request-id-pre", "auth-check-pre", "permission-check-pre",
				"permission-check-post", "auth-check-post", "request-id-post",
			},
		},
		{
			name: "request/response middleware pipeline",
			middlewares: []struct {
				pre  string
				post string
			}{
				{pre: "parse-body-pre", post: "parse-body-post"},
				{pre: "validate-pre", post: "validate-post"},
				{pre: "transform-pre", post: "transform-post"},
			},
			expectedExecutionOrder: []string{
				"parse-body-pre", "validate-pre", "transform-pre",
				"transform-post", "validate-post", "parse-body-post",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &mockWriter{}
			request := NewHTTPRequest()
			handler := func(writer HTTPWriter, request HTTPRequest) {}
			var middlewares []MiddlewareFunc
			var tracker []string
			for _, mw := range tt.middlewares {
				middlewares = append(middlewares, testMiddlewareWrapper(mw.pre, mw.post, &tracker))
			}
			res := ApplyMiddlewares(writer, request, middlewares, handler)
			res()

			if !slices.Equal(tt.expectedExecutionOrder, tracker) {
				t.Errorf("middleware execution order %s did not match the expected order %s", tracker, tt.expectedExecutionOrder)
			}
		})
	}
}

func testMiddlewareWrapper(
	pre, post string, tracker *[]string,
) MiddlewareFunc {
	return func(writer HTTPWriter, request HTTPRequest, next func()) {
		*tracker = append(*tracker, pre)
		next()
		*tracker = append(*tracker, post)
	}
}

type mockWriter struct{}

func (h *mockWriter) Response(payload string, statusCode int) {}

func (h *mockWriter) Header() Header {
	return &header{}
}
func (h *mockWriter) addHeader(header Header) {}

func TestGetMiddlewares(t *testing.T) {
	mw1 := func(writer HTTPWriter, request HTTPRequest, next func()) {}
	mw2 := func(writer HTTPWriter, request HTTPRequest, next func()) {}
	mw3 := func(writer HTTPWriter, request HTTPRequest, next func()) {}
	mw4 := func(writer HTTPWriter, request HTTPRequest, next func()) {}

	tests := []struct {
		name                string
		node                *node
		expectedMiddlewares []MiddlewareFunc
	}{
		{
			name: "node with no parent and no middlewares",
			node: &node{
				path:        "/",
				middlewares: []MiddlewareFunc{},
			},
			expectedMiddlewares: []MiddlewareFunc{},
		},
		{
			name: "node with middlewares but no parent",
			node: &node{
				path:        "/users",
				middlewares: []MiddlewareFunc{mw1, mw2},
			},
			expectedMiddlewares: []MiddlewareFunc{mw1, mw2},
		},
		{
			name: "node with parent that has middlewares",
			node: &node{
				path:        "/users/:id",
				middlewares: []MiddlewareFunc{mw3},
				parent: &node{
					path:        "/users",
					middlewares: []MiddlewareFunc{mw1, mw2},
					parent:      nil,
				},
			},
			expectedMiddlewares: []MiddlewareFunc{mw1, mw2, mw3},
		},
		{
			name: "nested nodes with multiple levels of parents",
			node: &node{
				path:        "/api/v1/users/:id/posts/:postId",
				middlewares: []MiddlewareFunc{mw4},
				parent: &node{
					path:        "/api/v1/users/:id",
					middlewares: []MiddlewareFunc{mw3},
					parent: &node{
						path:        "/api/v1",
						middlewares: []MiddlewareFunc{mw1, mw2},
						parent:      nil,
					},
				},
			},
			expectedMiddlewares: []MiddlewareFunc{mw1, mw2, mw3, mw4},
		},
		{
			name: "node with parent that has no middlewares",
			node: &node{
				path:        "/child",
				middlewares: []MiddlewareFunc{mw2},
				parent: &node{
					path:        "/parent",
					middlewares: []MiddlewareFunc{},
					parent:      nil,
				},
			},
			expectedMiddlewares: []MiddlewareFunc{mw2},
		},
		{
			name: "parent has middlewares, child has none",
			node: &node{
				path:        "/child",
				middlewares: []MiddlewareFunc{},
				parent: &node{
					path:        "/parent",
					middlewares: []MiddlewareFunc{mw1, mw2, mw3},
					parent:      nil,
				},
			},
			expectedMiddlewares: []MiddlewareFunc{mw1, mw2, mw3},
		},
		{
			name: "deep hierarchy of nodes with middlewares at different levels",
			node: &node{
				path:        "/level4",
				middlewares: []MiddlewareFunc{mw4},
				parent: &node{
					path:        "/level3",
					middlewares: []MiddlewareFunc{},
					parent: &node{
						path:        "/level2",
						middlewares: []MiddlewareFunc{mw2},
						parent: &node{
							path:        "/level1",
							middlewares: []MiddlewareFunc{mw1},
							parent:      nil,
						},
					},
				},
			},
			expectedMiddlewares: []MiddlewareFunc{mw1, mw2, mw4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middlewares := GetMiddlewares(tt.node)
			if len(middlewares) != len(tt.expectedMiddlewares) {
				t.Errorf("expected to get %v but got %v middlewares", len(tt.expectedMiddlewares), len(middlewares))
			}

			for i, mw := range middlewares {
				if reflect.ValueOf(mw).Pointer() != reflect.ValueOf(tt.expectedMiddlewares[i]).Pointer() {
					t.Errorf("expected middlewares to equal expected middlewares, but did not")
				}
			}
		})
	}
}
