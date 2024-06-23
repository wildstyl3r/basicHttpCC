package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buffer := make([]byte, 4096)
		_, err = conn.Read(buffer)

		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}

		req, err := readRequest(buffer)

		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(req.Url.Path)

		switch {
		case strings.HasPrefix(req.Url.Path, "/echo/"):
			conn.Write(
				newResponse(StatusOK).
					addHeader(ContentType, "text/plain").
					setBody(req.Url.Path[6:]).
					toBytes())
		case req.Url.Path == "/":
			conn.Write(newResponse(StatusOK).toBytes())
		case req.Url.Path == "/user-agent":
			conn.Write(
				newResponse(StatusOK).
					addHeader(ContentType, "text/plain").
					setBody(req.Headers[UserAgent]).
					toBytes())
		default:
			conn.Write(newResponse(StatusNotFound).toBytes())
		}
	}
}
