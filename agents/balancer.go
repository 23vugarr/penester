package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"masProject/agents/types"
	"net"
	"strings"
	"sync"
)

type Balancer struct {
	IP        string
	Port      string
	Agents    map[string]map[string]string
	tcpServer *types.TCPServer
}

func NewBalancer(ip, port string) *Balancer {
	return &Balancer{
		IP:     ip,
		Port:   port,
		Agents: make(map[string]map[string]string),
	}
}

func (b *Balancer) Run() {
	var wg sync.WaitGroup
	wg.Add(1)

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

	for {
		conn, err := b.tcpServer.Accept()
		if err != nil {
			fmt.Println("Error:", err)
		}
		log.Println("Accepting Connections...")

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
					ip := strings.Join(strings.Split(strings.Split(message.Content, ",")[0], ":")[1:], ":")
					maxLoad := strings.Split(strings.Split(message.Content, ",")[1], ":")[1]
					b.Agents[ip] = map[string]string{"ip": ip, "maxLoad": maxLoad}
					log.Println("Current agents: ", b.Agents)
				case "Instruction":
					log.Println("Message type: ", message.Type)
					go b.SendInstructions(wg)
				}
				log.Println("Message: ", message)
			}
		}(conn)

	}
}

func (b *Balancer) SendInstructions(wg *sync.WaitGroup) {
	defer wg.Done()
}
