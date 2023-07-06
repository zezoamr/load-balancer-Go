module github.com/zezoamr/load-balancer-Go

go 1.20

require (
    github.com/zezoamr/load-balancer-go/servers v0.0.0
    github.com/zezoamr/load-balancer-go/loadb v1.0.0
)

replace github.com/zezoamr/load-balancer-go/servers => ./servers
replace github.com/zezoamr/load-balancer-go/loadb => ./loadb

