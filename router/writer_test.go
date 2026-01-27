package router

import (
	"slices"
	"testing"
)

type mockConnection struct {
	written []byte
}

func (m *mockConnection) Write(p []byte) (int, error) {
	m.written = append(m.written, p...)
	return 0, nil
}

type mockHeader struct {
	key   HeaderType
	value string
}

func TestHttpWriter_FormatResponse(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		statusCode    int
		headers       []mockHeader
		expectedWrite []byte
		method        Request
	}{
		{
			name:       "GET request with 200 status and single header",
			method:     Get,
			payload:    "Hello World",
			statusCode: 200,
			headers: []mockHeader{
				{key: Host, value: "example.com"},
			},
			expectedWrite: []byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\nHost: example.com\r\n\r\nHello World"),
		},
		{
			name:       "POST request with 201 status and content length",
			method:     Post,
			payload:    `{"id":1,"name":"test"}`,
			statusCode: 201,
			headers: []mockHeader{
				{key: ContentType, value: "application/json"},
			},
			expectedWrite: []byte("HTTP/1.1 201 Created\r\nContent-Length: 22\r\nContent-Type: application/json\r\n\r\n{\"id\":1,\"name\":\"test\"}"),
		},
		{
			name:       "GET request with 404 status and no payload",
			method:     Get,
			payload:    "",
			statusCode: 404,
			headers: []mockHeader{
				{key: Host, value: "api.example.com"},
			},
			expectedWrite: []byte("HTTP/1.1 404 Not Found\r\nHost: api.example.com\r\n\r\n"),
		},
		{
			name:       "POST request with 500 error and multiple headers",
			method:     Post,
			payload:    `Internal Server Error`,
			statusCode: 500,
			headers: []mockHeader{
				{key: ContentType, value: "text/plain"},
				{key: Host, value: "error.example.com"},
			},
			expectedWrite: []byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 21\r\nContent-Type: text/plain\r\nHost: error.example.com\r\n\r\nInternal Server Error"),
		},
		{
			name:       "GET request with 301 redirect and location header",
			method:     Get,
			payload:    "",
			statusCode: 301,
			headers: []mockHeader{
				{key: Location, value: "/new-path"},
			},
			expectedWrite: []byte("HTTP/1.1 301 Moved Permanently\r\nLocation: /new-path\r\n\r\n"),
		},
		{
			name:       "POST request with 400 bad request",
			method:     Post,
			payload:    `{"error":"Invalid input"}`,
			statusCode: 400,
			headers: []mockHeader{
				{key: ContentType, value: "application/json"},
			},
			expectedWrite: []byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 25\r\nContent-Type: application/json\r\n\r\n{\"error\":\"Invalid input\"}"),
		},
		{
			name:       "GET request with 200 and authorization header",
			method:     Get,
			payload:    `{"token":"abc123"}`,
			statusCode: 200,
			headers: []mockHeader{
				{key: Authorization, value: "Bearer token123"},
				{key: ContentType, value: "application/json"},
			},
			expectedWrite: []byte("HTTP/1.1 200 OK\r\nContent-Length: 18\r\nAuthorization: Bearer token123\r\nContent-Type: application/json\r\n\r\n{\"token\":\"abc123\"}"),
		},
		{
			name:       "POST request with 204 no content",
			method:     Post,
			payload:    "",
			statusCode: 204,
			headers: []mockHeader{
				{key: Host, value: "api.example.com"},
			},
			expectedWrite: []byte("HTTP/1.1 204 No Content\r\nHost: api.example.com\r\n\r\n"),
		},
		{
			name:       "GET request with 403 forbidden",
			method:     Get,
			payload:    `Access Denied`,
			statusCode: 403,
			headers: []mockHeader{
				{key: ContentType, value: "text/plain"},
			},
			expectedWrite: []byte("HTTP/1.1 403 Forbidden\r\nContent-Length: 13\r\nContent-Type: text/plain\r\n\r\nAccess Denied"),
		},
		{
			name:       "POST request with complex JSON payload",
			method:     Post,
			payload:    `{"user":{"id":123,"email":"test@example.com","roles":["admin","user"]}}`,
			statusCode: 200,
			headers: []mockHeader{
				{key: ContentType, value: "application/json"},
				{key: Host, value: "api.example.com"},
			},
			expectedWrite: []byte("HTTP/1.1 200 OK\r\nContent-Length: 71\r\nContent-Type: application/json\r\nHost: api.example.com\r\n\r\n{\"user\":{\"id\":123,\"email\":\"test@example.com\",\"roles\":[\"admin\",\"user\"]}}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockConnection{}
			writer := NewHTTPWriter(mockConn, tt.method)
			for _, h := range tt.headers {
				writer.Header().Add(h.key, h.value)
			}
			writer.Response(tt.payload, tt.statusCode)

			if !slices.Equal(mockConn.written, tt.expectedWrite) {
				t.Errorf("expected write to be %s but got %s", tt.expectedWrite, mockConn.written)
			}
		})
	}
}
