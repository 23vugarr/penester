package penester

type Penester struct {
	BalancerIp     string
	instructionSet []string
}

func NewPenester(balancerIp string) *Penester {
	return &Penester{
		BalancerIp:     balancerIp,
		instructionSet: nil,
	}
}

func (p *Penester) ResolveInstructions(instructions []string) {
	p.instructionSet = instructions
}
