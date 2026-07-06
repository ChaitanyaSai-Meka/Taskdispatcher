package gatekeeper

type Gatekeeper struct {
	slots chan struct{}
}


func New(maxWorkers int) *Gatekeeper{
	if maxWorkers < 1 {
 		maxWorkers = 1
 	}
	return &Gatekeeper{
		slots:make(chan struct{},maxWorkers),
	}
}

func (g *Gatekeeper) TryAcquire()bool{
	select {
	case g.slots<-struct{}{}:
		return true
	default:
		return false
	}
	
}

func (g *Gatekeeper) Release(){
	<-g.slots
}
