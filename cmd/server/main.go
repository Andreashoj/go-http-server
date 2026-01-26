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
// Test get query param & get url param [X]
// Router route tests [X]
// Route grouping [X]
// Middlewares [X]
// Router config - allowed headers that is attached to all responses and same with middlewares [X]
// Middlewares test [X] []
// Create test for writer []
// Restructure for library package []

func main() {
	r := router.NewRouter()

	r.Use(func(writer router.HTTPWriter, request router.HTTPRequest, next func()) {
		fmt.Println("upper router middleware")
		writer.Header().Add(router.Host, "example.comz")
		next()
		fmt.Println("upper router middleware out")
	})

	r.Use(func(writer router.HTTPWriter, request router.HTTPRequest, next func()) {
		fmt.Println("upper router middleware 2")
		//writer.Header().Add(router.Host, "example.comz")
		next()
		fmt.Println("upper router middleware out 2")
	})

	r.Group("/user", func(u router.Router) {
		u.Use(func(writer router.HTTPWriter, request router.HTTPRequest, next func()) {
			fmt.Println("Yooo")
			next()
			fmt.Println("SUUUP")
		})

		u.Post("/my-other-route", func(writer router.HTTPWriter, request router.HTTPRequest) {
			fmt.Println("my other route!")
		})

		u.Get("/tester", func(writer router.HTTPWriter, request router.HTTPRequest) {
			writer.FormatResponse("YOOOOO BOI", 200)
		})

		u.Group("/anz", func(a router.Router) {

			a.Get("/yo", func(writer router.HTTPWriter, request router.HTTPRequest) {
				writer.FormatResponse("THIS BOY CRAZY", 200)
			})
		})
	})

	r.Get("/url/:id", func(writer router.HTTPWriter, request router.HTTPRequest) {
		// Retrieve id from query param ?
		urlParam, _ := request.GetURLParam("id")
		fmt.Println("param", urlParam)

		writer.FormatResponse(fmt.Sprintf("HERE YOU GO: %s", urlParam), 200)
	})

	r.Post("/url", func(w router.HTTPWriter, r router.HTTPRequest) {
		payload := r.Body()
		fmt.Println("payload", payload)

		param, err := r.GetQueryParam("tester")
		if err != nil {
			fmt.Printf("no query param with that name: %s", err)
		}
		fmt.Println(param)

		w.FormatResponse("epicposthasd fasdf asdf asd", 200)
	})

	if err := server.StartServer(":8080", r); err != nil {
		fmt.Printf("failed starting HTTP server: %s", err)
		return
	}

	// Makes requests to HTTP server with different kinds of http methods and bodies
	tests.CreateRequests(router.Post)

	select {}
}
