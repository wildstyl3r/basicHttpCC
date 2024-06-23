package main

import (
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	Url     *url.URL
	Headers map[string]string
	Body    string
}

func ReadRequest(b []byte) (Request, error) {
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

	headers := make(map[string]string)

	for _, headerLine := range head[1:] {
		if header, value, found := strings.Cut(headerLine, ": "); found {
			headers[header] = value
		}
	}

	body := r[bodyStart+len("\r\n\r\n") : strings.IndexByte(r, 0)]

	return Request{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
	}, nil
}
