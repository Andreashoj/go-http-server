package server

import (
	"fmt"
	"net"

	"github.com/Andreashoj/go-http-server/internal/parser"
	"github.com/Andreashoj/go-http-server/internal/router"
	"github.com/Andreashoj/go-http-server/internal/serializer"
)

func StartServer(port string, router router.Router) error {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		return fmt.Errorf("Failed creating listener for TCP on port: 8080, with error %s\n", err)
	}

	go func() {
		for {
			cn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Couldn't accept incoming TCP request with error: %s", err)
			}

			go func() {
				defer cn.Close()
				parser := parser.Listen(cn)
				request, err := parser.Parse()
				if err != nil {
					fmt.Printf("failed parsing http request: %s", err)
					return
				}

				route := router.FindMatchingRoute(request)

				// Writer
				writer := serializer.NewWriter(cn)
				route.Handler(writer)
			}()
		}
	}()

	return nil
}
