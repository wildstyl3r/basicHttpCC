package main

import "strconv"

type Response struct {
	Status  string
	Headers map[string]string
	Body    string
}

func newResponse(status string) *Response {
	return &Response{
		Status:  status,
		Headers: map[string]string{ContentLength: "0"},
	}
}

func (r *Response) setBody(body string) *Response {
	r.Body = body
	r.Headers[ContentLength] = strconv.Itoa(len(body))
	return r
}

func (r *Response) addHeader(header string, val string) *Response {
	r.Headers[header] = val
	return r
}

func (r *Response) toBytes() []byte {
	var headers string
	for header, val := range r.Headers {
		headers += header + ": " + val + "\r\n"
	}
	return []byte(
		"HTTP/1.1 " + r.Status + "\r\n" +
			headers + "\r\n" +
			r.Body)
}
