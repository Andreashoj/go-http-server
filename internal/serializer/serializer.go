package serializer

import (
	"fmt"
	"net"
)

type HTTPWriter interface {
	Writer(payload string)
}

type httpWriter struct {
	conn net.Conn
}

func NewWriter(conn net.Conn) HTTPWriter {
	return &httpWriter{
		conn: conn,
	}
}

func (h *httpWriter) Writer(payload string) {
	// Use conn
	fmt.Fprint(h.conn, payload)
}
