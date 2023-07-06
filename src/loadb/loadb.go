package loadb

import (
	"fmt"
	"net/http"

	"./servers"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []servers.Server
}

func newLoadBalancer(port string, servers []servers.Server) *LoadBalancer {

	return &LoadBalancer{
		port:            port,
		servers:         servers,
		roundRobinCount: 0,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() servers.Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.isAlive() {
		lb.roundRobinCount = lb.roundRobinCount + 1
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount = lb.roundRobinCount + 1
	//fmt.Println("roundrobin count", lb.roundRobinCount)
	return server
}
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Println("forwarding request to address: ", targetServer.getAddress())
	targetServer.serve(rw, r)
}
