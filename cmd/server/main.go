package main

import (
	"fmt"

	"github.com/Andreashoj/go-http-server/internal/router"
	"github.com/Andreashoj/go-http-server/internal/server"
	"github.com/Andreashoj/go-http-server/internal/tests"
)

func main() {
	r := router.NewRouter()

	r.Post("test", func(w router.HTTPWriter) {
		w.Header().Add(router.ContentLength, "10")
		w.Respond("epic post handler yo\n", 200)
		// Use own custom writer, that write will be used to format the request
	})

	if err := server.StartServer(":8080", r); err != nil {
		fmt.Printf("failed starting HTTP server: %s", err)
		return
	}

	// Makes requests to HTTP server with different kinds of http methods and bodies
	tests.CreateRequests(tests.Post)

	select {}
}
