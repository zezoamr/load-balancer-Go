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
	resp, err := http.Get("http://google.com/")
	handleError(err)
	return (resp.StatusCode == 200)
}

func (s *SimpleServer) serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.isAlive() {
		lb.roundRobinCount = lb.roundRobinCount + 1
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	//fmt.Println("roundrobin count", lb.roundRobinCount)
	return server
}
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Println("forwarding request to address: ", targetServer.getAddress())
	targetServer.serve(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.bing.com"),
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
