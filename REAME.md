# Penester - A Conceptual Multi-Agent Penetration Testing App

---
## Penester Architecture
1. Balancer [Controller] - Controls the agents that executes command sent by Balancer, receives conditions of agents and load balances the pipeline.
2. Agents - There can be mutliple agents and they are connected to balancer waiting for instructions from it, they are constantly sending heartbeats to balancer. 
3. Penester - API to accept pipelines to Balancer application, the rest is handled by Balancer application itself.

## How to run
Running balancer application if you have not compiled the application yet
```bash
go run .\balancer.go --host 127.0.0.1 --port 8999
```
Running agents if you have not compiled the application yet
```bash
go run .\agent.go --host 127.0.0.1 --port 9090 --maxload 3 --balancer 127.0.0.1:8999
go run .\agent.go --host 127.0.0.1 --port 9095 --maxload 5 --balancer 127.0.0.1:8999
```
Running the penester app if you have not compiled the application yet
```bash
go run .\penester.go --balancer 127.0.0.1:8999 --path .\test\example.yaml
```
## Example of pipeline yaml
```yaml
website: www.example.com

pipeline:
  portScan:
      start: 3000
      end: 50000
      type: tcp # udp
  dirScan:
      dirTxt: .\dir.txt

outputDir: .
```

## Future Development
A modern UI together with LLM support can be added for consultation of security vulnerabilities and for commenting on the results of pipeline. This application best suited for companies whom are constantly working on new projects and wants to automate initial steps of penetration testing with getting comments on result of pipeline from LLMs.