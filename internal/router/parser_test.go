package router

import (
	"bufio"
	"strings"
	"testing"
)

func TestParse_Startline(t *testing.T) {
	requests := []string{
		"GET / HTTP/1.1\r\n\r\n",
		"POST / HTTP/1.1\r\n\r\n",
		"DELETE / HTTP/1.1\r\n\r\n",
		"PUT / HTTP/1.1\r\n\r\n",
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
	requests := []string{
		"GET / HTTP/1.1\r\nContent-Length: 10\r\nContent-Type: json/application\r\nAccess-Control-Allow-Origin: http://example.com\r\n\r\n",
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
}

func TestParse_HeadersInvalid(t *testing.T) {
	requests := []string{
		"GET / \r\nHost: example.com\r\n\r\n",
		"GET / HTTP/2.5\r\nHost: example.com\r\n\r\n",
		"/ HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"INVALID / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"GET HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Length: 50\r\n\r\n",
		"POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Length: -10\r\n\r\n",
		"POST /api/data HTTP/1.1\r\nHost: example.com\r\nContent-Length: abc\r\n\r\n",
		"GET / HTTP/1.1\r\nUser-Agent: test\r\n\r\n",
		"GET / HTTP/1.1\r\nHost example.com\r\n\r\n",
	}

	for _, req := range requests {
		reader := bufio.NewReader(strings.NewReader(req))
		_, err := Parse(reader)

		if err == nil {
			t.Errorf("expected to fail with invalid headers: %s", err)
		}
	}
}
