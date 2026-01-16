package main

import (
	"fmt"

	"github.com/Andreashoj/go-http-server/internal/server"
	"github.com/Andreashoj/go-http-server/internal/testing"
)

func main() {
	if err := server.StartServer(); err != nil {
		fmt.Printf("failed starting HTTP server: %s", err)
		return
	}

	// Makes requests to HTTP server with different kinds of http methods and bodies
	testing.CreateTestRequests()

	select {}
}
