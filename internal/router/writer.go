package router

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HTTPWriter interface {
	Respond(payload string, statusCode int)
	Header() Header
	addHeader(header Header)
}

type httpWriter struct {
	conn    net.Conn
	method  Request
	headers []Header
}

func NewHTTPWriter(conn net.Conn, route *route) HTTPWriter {
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
	response.WriteString(fmt.Sprintf("HTTP/1.1 %s %s\r\n", status, h.method))

	// Headers
	if h.method == Post {
		response.WriteString(fmt.Sprintf("%s: %v\r\n", ContentLength, len(payload)))
	}

	for _, header := range h.headers {
		for key, value := range header.Get() {
			response.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}

	// Required empty line between body headers
	response.WriteString("\r\n")
	response.WriteString(payload)

	// Body
	fmt.Fprint(h.conn, response.String())
}

func (h *httpWriter) addHeader(header Header) {
	h.headers = append(h.headers, header)
}
