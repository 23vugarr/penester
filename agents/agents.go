package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"masProject/agents/types"
	"masProject/penester/pkg"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Agent struct {
	IP           string
	Port         string
	BalancerIp   string
	MaxLoad      int
	currentLoad  int
	tcpServer    *types.TCPServer
	currentTasks map[string]*Task
	doneTasks    map[string]*Task
}

type Task struct {
	state      bool
	trackingId string
	startTime  time.Time
	endTime    time.Time
}

func NewAgent(ip, port, balancerIp string, maxLoad int) *Agent {
	return &Agent{
		IP:           ip,
		Port:         port,
		BalancerIp:   balancerIp,
		MaxLoad:      maxLoad,
		currentTasks: make(map[string]*Task),
		doneTasks:    make(map[string]*Task),
	}
}

func (a *Agent) Run() {
	var wg sync.WaitGroup
	wg.Add(3)

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

	messageBytes, _ := json.Marshal(types.Message{Type: "Connect", Message: fmt.Sprintf("ip:%s%s, maxload:%d", a.IP, a.Port, a.MaxLoad)})

	if err := a.SendNotificationToBalancer(messageBytes); err != nil {
		log.Fatal(err)
	}
	log.Println("Balancer is informed...")

	go a.AcceptInstructions(&wg)
	go a.sendHeartbeat(&wg)
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
				err = a.ExecuteInstruction(message)
				if err != nil {
					return
				}
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

func (a *Agent) sendHeartbeat(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		time.Sleep(5 * time.Second)
		msg := types.Message{Type: "Heartbeat", Message: fmt.Sprintf("ip:%s%s, currentLoad:%d", a.IP, a.Port, a.currentLoad)}
		msgBytes, _ := json.Marshal(msg)
		err := a.SendNotificationToBalancer(msgBytes)
		if err != nil {
			return
		}
	}
}

func (a *Agent) ExecuteInstruction(message types.Message) error {
	trId := generateRandomId(5)
	a.addTask(trId)
	a.ExecutePortScanner(message, trId)
	//a.ExecuteDirScanner(message)
	return nil
}

func (a *Agent) ExecutePortScanner(message types.Message, trId string) {
	log.Println("Started port scanning...")
	doneCh := make(chan bool)
	fmt.Println(a.currentLoad)
	go pkg.PortScanner(message.Content.Pipeline.PortScan.Start, message.Content.Pipeline.PortScan.End, message.Content.Website, doneCh)

	done := <-doneCh

	if done {
		a.endTask(trId)
	}
}

func (a *Agent) ExecuteDirScanner(message types.Message) {
	log.Println("Started directory scanning...")
	pkg.DirScaner(message.Content.Website)
}

func (a *Agent) addTask(trackingId string) {
	a.currentLoad++
	task := Task{
		state:      true,
		trackingId: trackingId,
		startTime:  time.Now(),
	}
	a.currentTasks[trackingId] = &task
}

func (a *Agent) endTask(trackingId string) {
	a.currentLoad--
	for key, val := range a.currentTasks {
		if key == trackingId {
			delete(a.currentTasks, key)
			val.endTime = time.Now()
			a.doneTasks[key] = val
			log.Println("Done task: ", val)
			a.getInfo()
		}
		break
	}
}

func generateRandomId(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (a *Agent) getInfo() {
	log.Println("Current info...")
	log.Printf("Current load: %d, Done tasks: ", a.currentLoad)
	for _, val := range a.doneTasks {
		log.Printf("Tracking id: %s, start time: %s, end time: %s", val.trackingId, val.startTime.String(), val.endTime.String())
	}
}
