package testing

import (
	"bufio"
	"fmt"
	"net"
)

func CreateTestRequests() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("failed establishing connection to example.com through TCP")
		return
	}

	body := `{"name":"John"}`
	fmt.Fprintf(conn, "GET / HTTP/1.1\r\n"+
		"Host: example.com\r\n"+
		"Connection: close\r\n"+
		"Content-Length:%d\r\n\r\n"+
		"%s",
		len(body),
		body,
	)
	status, err := bufio.NewReader(conn).ReadString('\n')

	fmt.Println(status)
}
