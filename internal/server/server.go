package server

import (
	"fmt"
	"net"

	"github.com/Andreashoj/go-http-server/internal/router"
)

func StartServer(port string, r router.Router) error {
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
				parser := router.Listen(cn)
				request, err := parser.Parse()
				if err != nil {
					fmt.Printf("failed parsing http request: %s", err)
					return
				}

				route := r.FindMatchingRoute(request)
				if route == nil {
					return
				}

				// Writer
				writer := router.NewHTTPWriter(cn, route)
				route.Handler(writer, request)
			}()
		}
	}()

	return nil
}
