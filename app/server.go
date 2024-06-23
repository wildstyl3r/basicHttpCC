package main

import (
	"flag"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	workingDirectory := flag.String("directory", "/tmp/", "directory to expose files from")
	flag.Parse()

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	NewRouter(l).
		RegisterRouteWithArgs(Get, "/echo/$", func(r Request, arg ...string) *Response {
			return NewResponse(StatusOK).
				AddHeader(ContentType, "text/plain").
				SetBody(arg[0])
		}).
		RegisterRoute(Get, "/user-agent", func(r Request) *Response {
			return NewResponse(StatusOK).
				AddHeader(ContentType, "text/plain").
				SetBody(r.Headers[UserAgent])
		}).
		RegisterRoute(Get, "/", func(r Request) *Response {
			return NewResponse(StatusOK)
		}).
		RegisterRouteWithArgs(Get, "/files/$", func(r Request, arg ...string) *Response {
			file, err := os.Open(*workingDirectory + arg[0])
			if err != nil {
				return NewResponse(StatusNotFound)
			}

			var buffer [4096]byte
			n, err := file.Read(buffer[:])

			if err != nil {
				return NewResponse(StatusNotFound)
			}

			return NewResponse(StatusOK).
				AddHeader(ContentType, "application/octet-stream").
				SetBodyBinary(buffer[:n])

		}).
		Up()
}
