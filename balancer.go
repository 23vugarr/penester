package main

import "masProject/agents"

func main() {
	balancerIp := "127.0.0.1"
	balancerPort := ":8999"
	balancer := agents.NewBalancer(balancerIp, balancerPort)
	balancer.Run()
}
