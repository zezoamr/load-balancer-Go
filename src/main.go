package main

import (
	"fmt"
	"net/http"

	"./servers"
)

func main() {
	servers := []servers.Server{
		servers.Serverserver("https://www.bing.com"),
		servers.newSimpleServer("https://www.google.com"),
		servers.newSimpleServer("https://www.duckduckgo.com"),
	}
	lb := loadb.newLoadBalancer("8000", servers)
	handeRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}
	http.HandleFunc("/", handeRedirect)

	fmt.Println("serving requests at localhost/", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
