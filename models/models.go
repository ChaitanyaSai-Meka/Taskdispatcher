package models

type Task struct {
	ID int `json:"id"`
	TaskName string `json:"task_name"`
	Class int `json:"class"`
}