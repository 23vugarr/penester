package agents

type Agent struct {
	IP         string
	Port       string
	BalancerIp string
	MaxLoad    int
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

}
