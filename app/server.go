package main

import (
	"fmt"
	"net/url"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

type Request struct {
	Method  string
	Url     *url.URL
	Headers []string
	Body    string
}

func readRequest(b []byte) (Request, error) {
	r := string(b[:])
	bodyStart := strings.Index(r, "\r\n\r\n")
	if bodyStart == -1 {
		bodyStart = len(r) - 1
	}
	head := strings.Split(r[:bodyStart], "\r\n")
	requestLine := strings.Split(head[0], " ")
	method := requestLine[0]
	url, err := url.Parse(requestLine[1])
	if err != nil {
		return Request{}, err
	}

	body := r[bodyStart+len("\r\n\r\n"):]

	return Request{
		Method:  method,
		Url:     url,
		Headers: head[1:],
		Body:    body,
	}, nil
}

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

		if req.Url.Path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
}
