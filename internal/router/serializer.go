package router

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Andreashoj/go-http-server/internal/tests"
)

type HTTPWriter interface {
	Respond(payload string, statusCode int)
	Header() Header
	addHeader(header Header)
}

type httpWriter struct {
	conn    net.Conn
	method  tests.Request
	headers []Header
}

func NewWriter(conn net.Conn, route route) HTTPWriter {
	return &httpWriter{
		conn:   conn,
		method: route.Method,
	}
}

func (h *httpWriter) Respond(payload string, statusCode int) {
	// Create HTTP format response
	var response strings.Builder
	status := strconv.Itoa(statusCode)

	// Status line
	response.WriteString(fmt.Sprintf("HTTP/1.1 %s %s\r\n", h.method, status))

	// Headers
	for _, header := range h.headers {
		for key, value := range header.Get() {
			response.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}

	// Required empty line between body headers
	response.WriteString("\r\n\r\n")

	// These should be configurable through the router object
	// And also a writer.Header .. Should add to it

	// Body
	fmt.Println(response.String())

	fmt.Fprint(h.conn, payload)
}

func (h *httpWriter) addHeader(header Header) {
	h.headers = append(h.headers, header)
}
