package main

// cd src
// go run main.go server.go loadBalancer.go

import (
	"fmt"
	"net/http"
)

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
