package main

import (
	"flag"
	"log"
	"masProject/agents"
)

func main() {
	var agentIp string
	var agentPort string
	var balancerIp string
	var maxLoad int

	flag.StringVar(&agentIp, "host", "", "IP address of the agent")
	flag.StringVar(&agentPort, "port", "", "Port for the agent")
	flag.StringVar(&balancerIp, "balancer", "", "IP address of the balancer")
	flag.IntVar(&maxLoad, "maxload", 3, "Maximum load the agent can handle")
	flag.Parse()

	if agentIp == "" {
		log.Fatal("Please provide the agent IP address with --host argument")
	}
	if agentPort == "" {
		log.Fatal("Please provide the agent port with --port argument")
	}
	if balancerIp == "" {
		log.Fatal("Please provide the balancer IP with --balancer argument")
	}

	agentPort = ":" + agentPort

	agent := agents.NewAgent(agentIp, agentPort, balancerIp, maxLoad)
	agent.Run()
	// go run .\agent.go --host 127.0.0.1 --port 9090 --maxload 3 --balancer 127.0.0.1:8999
}
