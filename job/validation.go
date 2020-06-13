package job

var validators = []func(*Process) bool{
	func(p *Process) bool { return p.validateName() },
	func(p *Process) bool { return p.validateStartTime() },
	func(p *Process) bool { return p.validateStopTime() },
	func(p *Process) bool { return p.validateProcs() },
	func(p *Process) bool { return p.validateRetries() },
}

func alphaOnly(str string) bool {
	for _, c := range str {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}

func (p *Process) validateName() bool {
	nameLen := len(p.Name)
	if nameLen < 1 || nameLen > 32 || !alphaOnly(p.Name) {
		return false
	}
	return true
}

func (p *Process) validateStartTime() bool {
	return p.StartTime >= 0
}

func (p *Process) validateStopTime() bool {
	return p.StopTime >= 0
}

func (p *Process) validateProcs() bool {
	if p.Procs > 4 {
		return false
	}
	return true
}

func (p *Process) validateRetries() bool {
	if p.Retries > 4 {
		return false
	}
	return true
}
