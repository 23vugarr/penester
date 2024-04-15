package main

import (
	"flag"
	"log"
	"masProject/agents"
)

func main() {
	var balancerIp string
	var balancerPort string

	flag.StringVar(&balancerIp, "host", "", "IP address of the balancer")
	flag.StringVar(&balancerPort, "port", "", "Port for the balancer")
	flag.Parse()

	if balancerIp == "" {
		log.Fatal("Please provide the balancer IP address with --host argument")
	}
	if balancerPort == "" {
		log.Fatal("Please provide the balancer port with --port argument")
	}

	balancerPort = ":" + balancerPort

	balancer := agents.NewBalancer(balancerIp, balancerPort)
	balancer.Run()
}
