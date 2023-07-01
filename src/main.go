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

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func newSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &SimpleServer{
		addr:  addr,
		proxy: *httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func newLoadBalancer(port string, servers []Server) *LoadBalancer {

	return &LoadBalancer{
		port:            port,
		servers:         servers,
		roundRobinCount: 0,
	}
}

func (s *SimpleServer) getAddress() string {
	return s.addr
}

func (s *SimpleServer) isAlive() bool {
	return true
}

func (s *SimpleServer) serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	// use isalive ?
	var oldroundRobinCount = lb.roundRobinCount
	lb.roundRobinCount = (lb.roundRobinCount + 1) % len(lb.servers)
	return lb.servers[oldroundRobinCount]
}
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {

}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}
	lb := newLoadBalancer("8000", servers)
	handeRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}
	http.HandleFunc("/", handeRedirect)

	fmt.Println("serving requests at localhost/", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
