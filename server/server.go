package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	router2 "github.com/Andreashoj/go-http-server/router"
)

func StartServer(port string, r router2.Router) error {
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
				defer cn.Close() // Should be disabled if http keep-alive is set
				reader := bufio.NewReader(cn)
				request, err := router2.Parse(reader)
				if err != nil {
					if strings.Contains(err.Error(), "EOF") { // handles empty requests
						fmt.Println("empty request")
						return
					}

					fmt.Printf("failed parsing http request: %s", err)
					return
				}

				node, err := r.FindMatchingRoute(request)
				if err != nil {
					fmt.Printf("failed finding match for route: %s", err)
					return
				}

				if node.Route == nil {
					return
				}

				request.SetRouterURL(node.Route.Url)

				// Writer
				writer := router2.NewHTTPWriter(cn, node.Route.Method)
				middlewares := router2.GetMiddlewares(node)
				handler := router2.ApplyMiddlewares(writer, request, middlewares, node.Route.Handler)
				handler()
			}()
		}
	}()

	return nil
}
