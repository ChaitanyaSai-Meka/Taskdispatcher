package dispatcher

import (
	"sync"
	"time"

	"github.com/ChaitanyaSai-Meka/Taskdispatcher/internal/gatekeeper"
	"github.com/ChaitanyaSai-Meka/Taskdispatcher/internal/worker"
	"github.com/ChaitanyaSai-Meka/Taskdispatcher/models"
)

type Dispatcher struct {
	gk *gatekeeper.Gatekeeper

	newTaskCh chan models.Task
	doneCh    chan models.Task
	timeoutCh chan models.Task

	queue1 []models.Task
	queue2 []models.Task
	queue3 []models.Task

	resultsMu sync.RWMutex
	results   map[int]models.Task
}

func New(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		gk:        gatekeeper.New(maxWorkers),
		newTaskCh: make(chan models.Task),
		doneCh:    make(chan models.Task),
		timeoutCh: make(chan models.Task),
		results:   make(map[int]models.Task),
	}
}

func (d *Dispatcher) Submit(t models.Task) {
	d.newTaskCh <- t
}

func (d *Dispatcher) Status(id int) (models.Task, bool) {
	d.resultsMu.RLock()
	defer d.resultsMu.RUnlock()

	t, ok := d.results[id]
	return t, ok
}

func (d *Dispatcher) Run() {
	for {
		select {
		case t := <-d.newTaskCh:
			d.handleNewTask(t)

		case t := <-d.doneCh:
			d.gk.Release()
			d.recordResult(t)
			d.tryDispatchNext()

		case t := <-d.timeoutCh:
			d.handleTimeout(t)
		}
	}
}

func (d *Dispatcher) handleNewTask(t models.Task) {
	if d.gk.TryAcquire() {
		d.startWorker(t)
		return
	}

	t.Status = models.StatusQueued

	switch t.Class {
	case models.Class1:
		d.queue1 = append(d.queue1, t)
	case models.Class2:
		d.queue2 = append(d.queue2, t)
		time.AfterFunc(15*time.Second, func() { d.timeoutCh <- t })
	case models.Class3:
		d.queue3 = append(d.queue3, t)
		time.AfterFunc(5*time.Second, func() { d.timeoutCh <- t })
	default:
		t.Status = models.StatusFailed
	}

	d.recordResult(t)
}

func (d *Dispatcher) tryDispatchNext() {
	t, ok := d.popNextQueuedTask()
	if !ok {
		return
	}

	if !d.gk.TryAcquire() {
		d.pushFront(t)
		return
	}

	d.startWorker(t)
}

func (d *Dispatcher) handleTimeout(t models.Task) {
	if !d.removeQueuedTask(t) {
		return
	}

	t.Status = models.StatusBusy
	d.recordResult(t)
}

func (d *Dispatcher) startWorker(t models.Task) {
	t.Status = models.StatusRunning
	d.recordResult(t)
	go worker.Run(t, d.doneCh)
}

func (d *Dispatcher) recordResult(t models.Task) {
	d.resultsMu.Lock()
	defer d.resultsMu.Unlock()

	d.results[t.ID] = t
}

func (d *Dispatcher) popNextQueuedTask() (models.Task, bool) {
	if len(d.queue3) > 0 {
		return popFront(&d.queue3), true
	}
	if len(d.queue2) > 0 {
		return popFront(&d.queue2), true
	}
	if len(d.queue1) > 0 {
		return popFront(&d.queue1), true
	}

	return models.Task{}, false
}

func (d *Dispatcher) pushFront(t models.Task) {
	switch t.Class {
	case models.Class3:
		d.queue3 = append([]models.Task{t}, d.queue3...)
	case models.Class2:
		d.queue2 = append([]models.Task{t}, d.queue2...)
	default:
		d.queue1 = append([]models.Task{t}, d.queue1...)
	}
}

func (d *Dispatcher) removeQueuedTask(t models.Task) bool {
	switch t.Class {
	case models.Class2:
		return removeTaskByID(&d.queue2, t.ID)
	case models.Class3:
		return removeTaskByID(&d.queue3, t.ID)
	default:
		return false
	}
}

func popFront(queue *[]models.Task) models.Task {
	t := (*queue)[0]
	*queue = (*queue)[1:]
	return t
}

func removeTaskByID(queue *[]models.Task, id int) bool {
	for i, t := range *queue {
		if t.ID == id {
			*queue = append((*queue)[:i], (*queue)[i+1:]...)
			return true
		}
	}

	return false
}
