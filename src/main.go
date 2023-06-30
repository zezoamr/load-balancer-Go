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
	address string
	proxy   httputil.ReverseProxy
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func newSimpleServer(address string) *SimpleServer {
	serverUrl, err := url.Parse(address)
	handleError(err)

	return &SimpleServer{
		address: address,
		proxy:   *httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func newLoadBalancer(port string, servers []Server) *LoadBalancer {

	return &LoadBalancer{
		port:            port,
		servers:         servers,
		roundRobinCount: 0,
	}
}

func main() {
	fmt.Println("hello world")
}
