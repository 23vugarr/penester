package agents

import (
	"fmt"
	"log"
	"masProject/agents/types"
	"net"
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
	if err := a.SetTcpServer(); err != nil {
		log.Fatal(err)
	}
	defer func(tcpServer *types.TCPServer) {
		err := tcpServer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(a.tcpServer)

	log.Println("TCP server created...")

	conn, err := net.Dial("tcp", a.BalancerIp)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	fmt.Println("Connected to Balancer Server...")
	//reader := bufio.NewReader(os.Stdin)

	//message, _ := reader.ReadString('\n')
	message := fmt.Sprintf("ip addr: %s%s, maxload: %d", a.IP, a.Port, a.MaxLoad)

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := a.tcpServer.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go func(conn net.Conn) {
			log.Println("accepting connections....")

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

	//if err := a.SendNotificationToBalancer(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//go func() {
	//	err := a.AcceptInstructions()
	//	if err != nil {
	//
	//	}
	//}()
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

func (a *Agent) SendNotificationToBalancer() error {
	return nil
}

func (a *Agent) AcceptInstructions() error {
	return nil
}

func (a *Agent) ExecuteInstruction() error {
	return nil
}
