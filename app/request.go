package main

import (
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Url     *url.URL
	Headers map[string]string
	Body    []byte
}

func ReadRequest(b []byte) (Request, error) {
	r := string(b[:])
	headEnd := strings.Index(r, "\r\n\r\n")
	if headEnd == -1 {
		headEnd = len(r) - 1
	}

	head := strings.Split(r[:headEnd], "\r\n")
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

	bodyStart := headEnd + len("\r\n\r\n")

	var body []byte
	if lengthStr, present := headers[ContentLength]; present {
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return Request{}, err
		}
		body = b[bodyStart : bodyStart+length]
	} else {
		body = b[bodyStart : bodyStart+strings.IndexByte(r, 0)]
	}

	return Request{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
	}, nil
}
