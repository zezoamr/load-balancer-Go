package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func handleError(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

type Server interface {
	getAddress() string
	isAlive() bool
	serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	addr  string
	proxy httputil.ReverseProxy
}

func newSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &SimpleServer{
		addr:  addr,
		proxy: *httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func (s *SimpleServer) getAddress() string {
	return s.addr
}

func (s *SimpleServer) isAlive() bool {
	resp, err := http.Get(s.addr)
	handleError(err)
	return (resp.StatusCode == 200)
}

func (s *SimpleServer) serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}
