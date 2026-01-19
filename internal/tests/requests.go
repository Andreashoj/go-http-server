package tests

import (
	"bufio"
	"fmt"
	"net"

	router "github.com/Andreashoj/go-http-server/internal/router"
)

func CreateRequests(methods ...router.Request) {
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		fmt.Println("failed establishing connection to example.com through TCP")
		return
	}

	for _, method := range methods {
		switch method {
		case router.Post:
			createPostRequest(conn)
			break
		case router.Get:
			createGetRequest(conn)
			break
		case router.Put:
			createPutRequest(conn)
			break
		case router.Delete:
			createDeleteRequest(conn)
			break
		}
	}
}

func createPostRequest(conn net.Conn) {
	body := `{"name":"John"}`
	fmt.Fprintf(conn, "POST /url?tester=123 HTTP/1.1\r\n"+
		"Host: example.com\r\n"+
		"Connection: close\r\n"+
		"Content-Length:%d\r\n\r\n"+
		"%s",
		len(body),
		body,
	)

	_, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("failed reading response from HTTP server: %s\n", err)
	}
}

func createGetRequest(conn net.Conn) {
	fmt.Fprintf(conn, "GET / HTTP/1.1\r\n"+
		"Host: example.com\r\n"+
		"Connection: close\r\n\r\n",
	)

	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("failed reading response from HTTP server: %s", err)
	}
	fmt.Println(status)
}

func createPutRequest(conn net.Conn) {
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
	if err != nil {
		fmt.Printf("failed reading response from HTTP server: %s", err)
	}
	fmt.Println(status)
}

func createDeleteRequest(conn net.Conn) {
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
	if err != nil {
		fmt.Printf("failed reading response from HTTP server: %s", err)
	}
	fmt.Println(status)
}
