package penester

import (
	"encoding/json"
	"fmt"
	"masProject/agents/types"
	"masProject/penester/pkg"
	"net"
)

type Penester struct {
	BalancerIp string
	filePath   string
	Config     *pkg.Config
}

func NewPenester(balancerIp, filePath string) *Penester {
	return &Penester{
		BalancerIp: balancerIp,
		filePath:   filePath,
	}
}

func (p *Penester) ResolveInstructions() {
	p.Config = pkg.NewConfig(p.filePath)
	p.Config.LoadConfig()
}

func (p *Penester) SubmitToBalancer() error {
	conn, err := net.Dial("tcp", p.BalancerIp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error closing the connection:", err)
		}
	}()
	messageBytes, _ := json.Marshal(types.Message{Type: "Submission", Message: "Instructions", Content: *p.Config})

	_, err = conn.Write(messageBytes)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
