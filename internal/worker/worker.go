package worker

import (
	"time"

	"github.com/ChaitanyaSai-Meka/Taskdispatcher/models"
)

func Run(task models.Task, doneCh chan<- models.Task) {
	duration := DurationFor(task.Class)
	time.Sleep(duration)

	task.Status = models.StatusDone

	doneCh <- task
}

func DurationFor(class models.JobClass) time.Duration {
	switch class {
	case models.Class1:
		return 5 * time.Second
	case models.Class2:
		return 10 * time.Second
	case models.Class3:
		return 15 * time.Second
	default:
		return 5 * time.Second
	}
}
