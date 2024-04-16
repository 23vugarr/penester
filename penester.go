package main

import (
	"flag"
	"masProject/penester"
)

func main() {
	var balancer string
	var path string

	flag.StringVar(&balancer, "balancer", "", "IP address of the balancer")
	flag.StringVar(&path, "path", "", "Port for the balancer")
	flag.Parse()

	penesterAgent := penester.NewPenester(balancer, path)
	penesterAgent.ResolveInstructions()
	err := penesterAgent.SubmitToBalancer()
	if err != nil {
		return
	}
	// go run .\penester.go --balancer 127.0.0.1:8999 --path .\test\example.yaml
}
