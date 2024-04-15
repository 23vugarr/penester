package agents

import (
	"fmt"
	"log"
	"masProject/agents/types"
	"net"
)

type Balancer struct {
	IP        string
	Port      string
	Agents    map[string]string
	tcpServer *types.TCPServer
}

func NewBalancer(ip, port string) *Balancer {
	return &Balancer{
		IP:     ip,
		Port:   port,
		Agents: nil,
	}
}

func (b *Balancer) Run() {
	if err := b.SetTcpServer(); err != nil {
		log.Fatal(err)
	}

	defer func(tcpServer *types.TCPServer) {
		err := tcpServer.Close()
		if err != nil {

		}
	}(b.tcpServer)

	log.Println("TCP server created...")

	for {
		conn, err := b.tcpServer.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go func(conn net.Conn) {
			log.Println("accepting connections....")
			//b.Agents[conn.RemoteAddr().String()] = "hello"
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

				fmt.Printf("Received: %s\n", buffer[:n])
			}
		}(conn)
	}

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

func (b *Balancer) GetAgents() {

}
