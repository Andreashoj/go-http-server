package router

import (
	"bufio"
	"strings"
	"testing"
)

func TestParse_Startline(t *testing.T) {
	requests := []string{
		"GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"POST / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"DELETE / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"PUT / HTTP/1.1\r\nHost: example.com\r\n\r\n",
	}

	for _, req := range requests {
		reader := bufio.NewReader(strings.NewReader(req))

		_, err := Parse(reader)

		if err != nil {
			t.Errorf("failed parsing request: %s", err)
		}
	}
}

func TestParse_StartlineShouldFail(t *testing.T) {
	requests := []string{
		"GET / HTTP/1.1\r\n",
		"POST / HTTP/1.1r\n\r\n",
		"DELETE HTTP/1.1\r\n\r\n",
		"/ HTTP/1.1\r\n\r\n",
	}

	for _, req := range requests {
		reader := bufio.NewReader(strings.NewReader(req))
		_, err := Parse(reader)

		if err == nil {
			t.Errorf("expected parser to fail: %s", err)
		}
	}
}

func TestParse_Headers(t *testing.T) {
	t.Run("valid headers", func(t *testing.T) {
		requests := []string{
			"POST /api/users HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: application/json\r\nContent-Length: 27\r\nAuthorization: Bearer token123\r\n\r\n",
			"GET /search?q=golang HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html\r\nCookie: session=abc123\r\n\r\n",
			"PUT /api/posts/42 HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 50\r\nAuthorization: Bearer token456\r\n\r\n",
			"DELETE /api/items/99 HTTP/1.1\r\nHost: example.com\r\nAuthorization: Bearer token789\r\nAccept: application/json\r\n\r\n",
			"PATCH /api/config HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: application/json\r\nContent-Length: 35\r\nETag: \"33a64df551\"\r\n\r\n",
			"HEAD / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			"OPTIONS /api/data HTTP/1.1\r\nHost: example.com\r\nAccess-Control-Request-Method: POST\r\nAccess-Control-Request-Headers: Content-Type\r\n\r\n",
		}

		for _, req := range requests {
			reader := bufio.NewReader(strings.NewReader(req))
			_, err := Parse(reader)

			if err != nil {
				t.Errorf("failed parsing request: %s", err)
			}
		}
	})

	t.Run("invalid headers", func(t *testing.T) {
		requests := []string{
			"POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Length: -10\r\n\r\n",
			"POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Length: abc\r\n\r\n",
			"GET / HTTP/1.1\r\nUser-Agent: test\r\n\r\n",
			"GET / HTTP/1.1\r\nHost example.com\r\n\r\n",
		}

		for _, req := range requests {
			reader := bufio.NewReader(strings.NewReader(req))
			_, err := Parse(reader)

			if err == nil {
				t.Errorf("expected to fail with invalid headers on: %s", req)
			}
		}
	})
}

// Add edge case for empty params
func TestParse_Params(t *testing.T) {
	requests := []struct {
		url   string
		key   string
		value string
	}{
		{
			url:   "GET /search?q=golang HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html\r\nCookie: session=abc123\r\n\r\n",
			key:   "q",
			value: "golang",
		},
		{
			url:   "GET /search?q=golang&sort=stars HTTP/1.1\r\nHost: example.com\r\n\r\n",
			key:   "sort",
			value: "stars",
		},
		{
			url:   "GET /search?q=hello+world HTTP/1.1\r\nHost: example.com\r\n\r\n",
			key:   "q",
			value: "hello world",
		},
		{
			url:   "GET /api?page=5&limit=10 HTTP/1.1\r\nHost: example.com\r\n\r\n",
			key:   "page",
			value: "5",
		},
		{
			url:   "GET /path?foo=bar HTTP/1.1\r\nHost: example.com\r\n\r\n",
			key:   "foo",
			value: "bar",
		},
	}

	for _, tt := range requests {
		reader := bufio.NewReader(strings.NewReader(tt.url))
		httpReq, err := Parse(reader)

		if err != nil {
			t.Errorf("failed, didn't expect there to be errors while parsing request: %s", err)
		}

		key, err := httpReq.GetQueryParam(tt.key)
		if err != nil {
			t.Errorf("failed getting query param: %s", err)
		}

		if key != tt.value {
			t.Errorf("failed, expected %s, but got %s", tt.value, key)
		}
	}
}

func TestParse_GetMethod(t *testing.T) {
	requests := []struct {
		request        string
		expectedMethod Request
	}{
		{
			request:        "GET / HTTP/1.1\r\nHost: api.example.com\r\nContent-Length: 10\r\nContent-Type: json/application\r\nAccess-Control-Allow-Origin: http://example.com\r\n\r\n",
			expectedMethod: Get,
		},
		{
			request:        "POST /api/users HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: application/json\r\nContent-Length: 27\r\nAuthorization: Bearer token123\r\n\r\n",
			expectedMethod: Post,
		},
		{
			request:        "GET /search?q=golang HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html\r\nCookie: session=abc123\r\n\r\n",
			expectedMethod: Get,
		},
		{
			request:        "PUT /api/posts/42 HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 50\r\nAuthorization: Bearer token456\r\n\r\n",
			expectedMethod: Put,
		},
		{
			request:        "DELETE /api/items/99 HTTP/1.1\r\nHost: example.com\r\nAuthorization: Bearer token789\r\nAccept: application/json\r\n\r\n",
			expectedMethod: Delete,
		},
		{
			request:        "PATCH /api/config HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: application/json\r\nContent-Length: 35\r\nETag: \"33a64df551\"\r\n\r\n",
			expectedMethod: Patch,
		},
		{
			request:        "HEAD / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			expectedMethod: Head,
		},
		{
			request:        "OPTIONS /api/data HTTP/1.1\r\nHost: example.com\r\nAccess-Control-Request-Method: POST\r\nAccess-Control-Request-Headers: Content-Type\r\n\r\n",
			expectedMethod: Options,
		},
	}

	for _, tt := range requests {
		reader := bufio.NewReader(strings.NewReader(tt.request))
		httpReq, err := Parse(reader)

		if err != nil {
			t.Errorf("test failed: %s", err)
		}

		if httpReq.Method() != tt.expectedMethod {
			t.Errorf("failed, expected %s to equal %s", httpReq.Method(), tt.expectedMethod)
		}
	}
}

func TestParse_ContentLength(t *testing.T) {
	requests := []struct {
		request               string
		expectedContentLength int
	}{
		{
			request:               "OPTIONS /api/data HTTP/1.1\r\nHost: example.com\r\nAccess-Control-Request-Method: POST\r\nAccess-Control-Request-Headers: Content-Type\r\n\r\n",
			expectedContentLength: 0,
		},
		{
			request:               "POST /api/users HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 24\r\n\r\n{\"name\":\"John\",\"age\":30}",
			expectedContentLength: 24,
		}, {
			request:               "POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Type: text/plain\r\nContent-Length: 11\r\n\r\nHello World",
			expectedContentLength: 11,
		},
		{
			request:               "PUT /api/config HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 32\r\n\r\n{\"setting\":\"value\",\"enabled\":true}",
			expectedContentLength: 32,
		},
		{
			request:               "PATCH /api/item HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 18\r\n\r\n{\"status\":\"active\"}",
			expectedContentLength: 18,
		},
	}

	for _, tt := range requests {
		reader := bufio.NewReader(strings.NewReader(tt.request))
		httpReq, err := Parse(reader)

		if err != nil {
			t.Errorf("test failed: %s", err)
		}

		if len(httpReq.Body()) != tt.expectedContentLength {
			t.Errorf("expected body length to be %v but got %v", tt.expectedContentLength, len(httpReq.Body()))
		}
	}
}

func TestParse_Body(t *testing.T) {
	requests := []struct {
		request string
		body    string
	}{
		{
			request: "PUT /api/config HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 34\r\n\r\n{\"setting\":\"value\",\"enabled\":true}",
			body:    "{\"setting\":\"value\",\"enabled\":true}",
		},
		{
			request: "POST /api/users HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 24\r\n\r\n{\"name\":\"John\",\"age\":30}",
			body:    "{\"name\":\"John\",\"age\":30}",
		},
		{
			request: "POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Type: text/plain\r\nContent-Length: 11\r\n\r\nHello World",
			body:    "Hello World",
		},
		{
			request: "PATCH /api/item HTTP/1.1\r\nHost: example.com\r\nContent-Type: application/json\r\nContent-Length: 19\r\n\r\n{\"status\":\"active\"}",
			body:    "{\"status\":\"active\"}",
		},
		{
			request: "POST /api/message HTTP/1.1\r\nHost: example.com\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nTest message!",
			body:    "Test message!",
		},
	}

	for _, tt := range requests {
		reader := bufio.NewReader(strings.NewReader(tt.request))
		httpReq, err := Parse(reader)

		if err != nil {
			t.Errorf("test failed: %s", err)
		}

		if httpReq.Body() != tt.body {
			t.Errorf("expected body to equal %s but got %s", tt.body, httpReq.Body())
		}
	}
}
