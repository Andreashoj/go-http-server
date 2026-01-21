package main

import (
	"fmt"

	"github.com/Andreashoj/go-http-server/internal/router"
	"github.com/Andreashoj/go-http-server/internal/server"
	"github.com/Andreashoj/go-http-server/internal/tests"
)

// TODO:
// Write tests for Parser [X]
// FEATURES: [X]
// Retrieve ID from URL [X]
// Test get query param & get url param
// Route grouping
// Middlewares
// Router config - allowed headers that is attached to all responses
// Create tests for the above if needed
// That concludes this project

func main() {
	r := router.NewRouter()

	r.Get("/url/:id", func(writer router.HTTPWriter, request router.HTTPRequest) {
		// Retrieve id from query param ?
		urlParam, _ := request.GetURLParam("id")
		fmt.Println("param", urlParam)

		writer.Respond(fmt.Sprintf("HERE YOU GO: %s", urlParam), 200)
	})

	r.Post("/url", func(w router.HTTPWriter, r router.HTTPRequest) {
		payload := r.Body()
		fmt.Println("payload", payload)

		param, err := r.GetQueryParam("tester")
		if err != nil {
			fmt.Printf("no query param with that name: %s", err)
		}
		fmt.Println(param)

		w.Respond("epicposthasd fasdf asdf asd", 200)
	})

	if err := server.StartServer(":8080", r); err != nil {
		fmt.Printf("failed starting HTTP server: %s", err)
		return
	}

	// Makes requests to HTTP server with different kinds of http methods and bodies
	tests.CreateRequests(router.Post)

	select {}
}
