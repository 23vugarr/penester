package main

import "masProject/agents"

func main() {
	agentIp := "127.0.0.1"
	agentPort := ":9090"
	balancerIp := "127.0.0.1:8999"
	maxLoad := 3
	agent := agents.NewAgent(agentIp, agentPort, balancerIp, maxLoad)
	agent.Run()
}
