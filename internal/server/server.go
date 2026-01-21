package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

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
				defer cn.Close() // Should be disabled if http keep-alive is set
				reader := bufio.NewReader(cn)
				request, err := router.Parse(reader)
				if err != nil {
					if strings.Contains(err.Error(), "EOF") { // handles empty requests
						fmt.Println("empty request")
						return
					}

					fmt.Printf("failed parsing http request: %s", err)
					return
				}

				route := r.FindMatchingRoute(request)
				request.SetRouterURL(route.Url)
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
