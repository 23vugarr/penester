package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"masProject/agents/types"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Balancer struct {
	IP        string
	Port      string
	Agents    map[string]*AgentInfo
	tcpServer *types.TCPServer
}

type AgentInfo struct {
	IP          string
	MaxLoad     int
	CurrentLoad int
	heartbeat   time.Time
}

func NewBalancer(ip, port string) *Balancer {
	return &Balancer{
		IP:     ip,
		Port:   port,
		Agents: make(map[string]*AgentInfo),
	}
}

func (b *Balancer) Run() {
	var wg sync.WaitGroup
	wg.Add(2)

	if err := b.SetTcpServer(); err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server is created...")
	defer func(tcpServer *types.TCPServer) {
		err := tcpServer.Close()
		if err != nil {

		}
	}(b.tcpServer)

	go b.GetMessages(&wg)
	go b.checkAgentHeartbeats(&wg)

	wg.Wait()

}

func (b *Balancer) SetTcpServer() error {
	server := types.NewTcpServer(b.IP, b.Port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	b.tcpServer = server

	return nil
}

func (b *Balancer) GetMessages(wg *sync.WaitGroup) {
	defer wg.Done()
	var mu sync.Mutex

	log.Println("Accepting Connections...")

	for {
		conn, err := b.tcpServer.Accept()
		if err != nil {
			fmt.Println("Error:", err)
		}

		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)
			buffer := make([]byte, 1024)

			for {
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				var message types.Message
				if err := json.Unmarshal(buffer[:n], &message); err != nil {
					log.Fatal(err)
					return
				}
				switch message.Type {
				case "Connect":
					ip := strings.Join(strings.Split(strings.Split(message.Message, ",")[0], ":")[1:], ":")
					maxLoad := strings.Split(strings.Split(message.Message, ",")[1], ":")[1]
					maxLoadInt, _ := strconv.Atoi(maxLoad)
					b.Agents[ip] = &AgentInfo{IP: ip, MaxLoad: maxLoadInt, CurrentLoad: 0, heartbeat: time.Now()}
					log.Printf("Agent connected: %s with max load %d\n", ip, maxLoadInt)
				case "Heartbeat":
					ip := strings.Join(strings.Split(strings.Split(message.Message, ",")[0], ":")[1:], ":")
					currLoad, _ := strconv.Atoi(strings.Join(strings.Split(strings.Split(message.Message, ",")[1], ":")[1:], ":"))
					b.getHeartbeats(&mu, ip, currLoad)
				case "Submission":
					log.Println("Got new task submission...")
					b.distributeTask(message)
				case "TaskDone":
					fmt.Println("will be done")
				}
			}
		}(conn)

	}
}

func (b *Balancer) distributeTask(message types.Message) {
	agent := b.selectLeastLoadedAgent()
	if agent == nil {
		log.Println("No agents available")
		return
	}
	go b.sendInstructions(agent, message)
}

func (b *Balancer) selectLeastLoadedAgent() *AgentInfo {
	var selected *AgentInfo
	for _, agent := range b.Agents {
		if selected == nil || float32(agent.CurrentLoad+1)/float32(agent.MaxLoad) < float32(selected.CurrentLoad+1)/float32(selected.MaxLoad) {
			selected = agent
		}
		fmt.Println(agent)
	}
	return selected
}

func (b *Balancer) sendInstructions(agent *AgentInfo, message types.Message) {
	agent.CurrentLoad++
	conn, err := net.Dial("tcp", agent.IP)
	if err != nil {
		log.Println("Error dialing agent:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	msg, _ := json.Marshal(message)
	if _, err := conn.Write(msg); err != nil {
		log.Println("Error sending instructions:", err)
	} else {
		log.Printf("Task sent to %s...", agent.IP)
	}
}

func (b *Balancer) checkAgentHeartbeats(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		for key, value := range b.Agents {
			if diff := time.Now().Sub(value.heartbeat); diff.Seconds() > 6 {
				fmt.Println(diff)
				delete(b.Agents, key)
				log.Println("Problem occured while getting heartbeat of", key)
				log.Println("Current agents: ", b.Agents)
			}
		}
		if len(b.Agents) == 0 {
			log.Println("No Agent in the Agent List...")
		}

		time.Sleep(10 * time.Second)
	}
}

func (b *Balancer) getHeartbeats(mu *sync.Mutex, id string, currLoad int) {
	mu.Lock()
	defer mu.Unlock()

	agent, exists := b.Agents[id]
	if !exists {
		log.Printf("Agent with ID %s not found", id)
		return
	}

	agent.heartbeat = time.Now()
	agent.CurrentLoad = currLoad
	log.Printf("Updated heartbeat of %s with current load of %d", id, currLoad)
}
