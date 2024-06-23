package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Router struct {
	routesWArgs map[string]map[string]func(Request, ...string) *Response
	routes      map[string]map[string]func(Request) *Response
	listener    net.Listener
}

func NewRouter(l net.Listener) *Router {
	return &Router{
		make(map[string]map[string]func(Request, ...string) *Response),
		make(map[string]map[string]func(Request) *Response),
		l,
	}
}

func (r *Router) RegisterRoute(method, path string, handler func(Request) *Response) *Router {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]func(Request) *Response)
	}

	r.routes[method][path] = handler
	return r
}

func (r *Router) RegisterRouteWithArgs(method, path string, handler func(Request, ...string) *Response) *Router {
	if r.routesWArgs[method] == nil {
		r.routesWArgs[method] = make(map[string]func(Request, ...string) *Response)
	}

	r.routesWArgs[method][path] = handler
	return r
}

func (r *Router) Up() {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		buffer := make([]byte, 4096)
		_, err = conn.Read(buffer)

		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}

		req, err := ReadRequest(buffer)

		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(req.Url.Path)

		gzip := false
		if strings.Contains(req.Headers[AcceptEncoding], "gzip") {
			gzip = true
		}

		response := NewResponse(StatusNotFound)

		if r.routesWArgs[req.Method] != nil {
			for route, handler := range r.routesWArgs[req.Method] {
				routeArgs := strings.Split(route, "$")
				if strings.HasPrefix(req.Url.Path, routeArgs[0]) {
					wip := req.Url.Path[len(routeArgs[0]):]
					var args []string
					for _, e := range routeArgs[1:] {
						if e != "" {
							next := strings.Index(wip, e)
							args = append(args, wip[0:next])
							wip = wip[next+len(e):]
						} else {
							args = append(args, wip)
						}
					}
					response = handler(req, args...)
					break
				}
			}
		}

		if r.routes[req.Method] != nil {
			for route, handler := range r.routes[req.Method] {
				if req.Url.Path == route {
					response = handler(req)
					break
				}
			}
		}
		if gzip && len(response.Body) > 0 {
			response.AddHeader(ContentEncoding, "gzip")
		}
		conn.Write(response.toBytes())
	}
}
