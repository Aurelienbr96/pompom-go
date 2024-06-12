package dto

import "time"

type TaskToCreate struct {
	TagId []int
	Task
}
type Task struct {
	Name        string
	Description string
	Duration    int64
	Date        time.Time
	UserId      int
}
