package router

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPWriter interface {
	Response(payload string, statusCode int)
	Header() Header
	addHeader(header Header)
}

type httpWriter struct {
	conn    Connection
	method  Request
	headers []Header
}

type Connection interface {
	Write(p []byte) (n int, err error)
}

func NewHTTPWriter(conn Connection, method Request) HTTPWriter {
	return &httpWriter{
		conn:   conn,
		method: method,
	}
}

func (h *httpWriter) Response(payload string, statusCode int) {
	// Create HTTP format response
	var response strings.Builder
	status := strconv.Itoa(statusCode)

	// Status line
	response.WriteString(fmt.Sprintf("HTTP/1.1 %s %s\r\n", status, getStatusMessage(statusCode)))

	// Headers
	if len(payload) > 0 {
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

func getStatusMessage(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 301:
		return "Moved Permanently"
	case 302:
		return "Found"
	case 304:
		return "Not Modified"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	default:
		return "Unknown"
	}
}
