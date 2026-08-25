package task

import "time"

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	Id          int
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
