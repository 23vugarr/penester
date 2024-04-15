package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"masProject/agents/types"
	"net"
	"sync"
)

type Agent struct {
	IP         string
	Port       string
	BalancerIp string
	MaxLoad    int
	tcpServer  *types.TCPServer
}

func NewAgent(ip, port, balancerIp string, maxLoad int) *Agent {
	return &Agent{
		IP:         ip,
		Port:       port,
		BalancerIp: balancerIp,
		MaxLoad:    maxLoad,
	}
}

func (a *Agent) Run() {
	var wg sync.WaitGroup
	wg.Add(2)

	if err := a.SetTcpServer(); err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server is created...")
	defer func(tcpServer *types.TCPServer) {
		err := tcpServer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(a.tcpServer)

	messageBytes, _ := json.Marshal(types.Message{Type: "Connect", Content: fmt.Sprintf("ip addr: %s%s, maxload: %d", a.IP, a.Port, a.MaxLoad)})

	if err := a.SendNotificationToBalancer(messageBytes); err != nil {
		log.Fatal(err)
	}
	log.Println("Balancer is informed...")

	go a.AcceptInstructions(&wg)
	wg.Wait()
}

func (a *Agent) SetTcpServer() error {
	a.tcpServer = types.NewTcpServer(a.IP, a.Port)
	err := a.tcpServer.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (a *Agent) AcceptInstructions(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		conn, err := a.tcpServer.Accept()
		if err != nil {
			fmt.Println("Error:", err)
		}
		log.Println("Accepting Instructions...")

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
				log.Println(message)
			}
		}(conn)

	}
}

func (a *Agent) SendNotificationToBalancer(message []byte) error {
	conn, err := net.Dial("tcp", a.BalancerIp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error closing the connection:", err)
		}
	}()

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (a *Agent) ExecuteInstruction() error {
	return nil
}
