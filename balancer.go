package main

import "masProject/agents"

func main() {
	balancer := agents.NewBalancer()
	balancer.Run()
}
