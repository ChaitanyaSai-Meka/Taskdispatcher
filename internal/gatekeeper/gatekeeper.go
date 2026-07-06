package gatekeeper

type Gatekeeper struct {
	slots chan struct{}
}


func New(maxWorkers int) *Gatekeeper{
	return &Gatekeeper{
		slots:make(chan struct{},maxWorkers),
	}
}

func (g *Gatekeeper) TryAquire()bool{
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
