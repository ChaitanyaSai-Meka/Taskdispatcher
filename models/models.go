package models

type JobClass int

const (
	Class1 JobClass = iota + 1
	Class2
	Class3
)

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
	StatusBusy    Status = "worker_busy"
)

type Task struct {
	ID       int      `json:"id"`
	TaskName string   `json:"task_name"`
	Class    JobClass `json:"class"`
	Status   Status   `json:"status"`
}
