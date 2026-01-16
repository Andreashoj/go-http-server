package server

import (
	"fmt"
	"net"

	"github.com/Andreashoj/go-http-server/internal/parser"
)

func StartServer() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("Failed creating listener for TCP on port: 8080, with error %s\n", err)
	}

	go func() {
		for {
			cn, err := ln.Accept()
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
				fmt.Println("received: ", request)
				response := "YOU LITTLE PRICK!"
				fmt.Fprint(cn, response)
			}()
		}
	}()

	return nil
}
