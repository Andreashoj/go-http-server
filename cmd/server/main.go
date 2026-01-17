package main

import (
	"fmt"
	"net"

	"github.com/Andreashoj/go-http-server/internal/router"
	"github.com/Andreashoj/go-http-server/internal/server"
	"github.com/Andreashoj/go-http-server/internal/testing"
)

func main() {
	r := router.NewRouter()

	r.Post("test", func(cn net.Conn) {
		fmt.Fprint(cn, "epic post handler")
	})

	if err := server.StartServer(":8080", r); err != nil {
		fmt.Printf("failed starting HTTP server: %s", err)
		return
	}

	// Makes requests to HTTP server with different kinds of http methods and bodies
	testing.CreateTestRequests()

	select {}
}
