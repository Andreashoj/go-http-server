package tests

import (
	"bufio"
	"fmt"
	"net"
)

type Request string

const (
	Post   Request = "POST"
	Get    Request = "GET"
	Put    Request = "PUT"
	Delete Request = "DELETE"
)

func CreateRequests(methods ...Request) {
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		fmt.Println("failed establishing connection to example.com through TCP")
		return
	}

	for _, method := range methods {
		switch method {
		case Post:
			createPostRequest(conn)
			break
		case Get:
			createGetRequest(conn)
			break
		case Put:
			createPutRequest(conn)
			break
		case Delete:
			createDeleteRequest(conn)
			break
		}
	}
}

func createPostRequest(conn net.Conn) {
	body := `{"name":"John"}`
	fmt.Fprintf(conn, "POST / HTTP/1.1\r\n"+
		"Host: example.com\r\n"+
		"Connection: close\r\n"+
		"Content-Length:%d\r\n\r\n"+
		"%s",
		len(body),
		body,
	)

	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("failed reading response from HTTP server: %s\n", err)
	}
	fmt.Println(status)
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
