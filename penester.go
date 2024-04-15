package main

import "masProject/penester"

func main() {
	balancerIp := "127.0.0.1:8999"
	instructions := make([]string, 0)
	penesterAgent := penester.NewPenester(balancerIp)
	penesterAgent.ResolveInstructions(instructions)
}
