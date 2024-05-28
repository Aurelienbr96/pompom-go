package model

import (
	taskDTO "pompom/go/src/dto"
	"time"
)

type TagInTask struct {
	Color string `db:"color" json:"color"`
}

type Task struct {
	ID int64 `db:"id" json:"id"`
	TaskToCreate
}

type TaskToCreate struct {
	Name        string    `db:"name" json:"name"`
	Date        time.Time `db:"date" json:"date"`
	Duration    int64     `db:"duration" json:"duration"`
	Description string    `db:"description" json:"description,omitempty"`
	TagId       int64     `db:"tagid" json:"tagid"`
}

type ExtendedTask struct {
	Task
	TagInTask
}

type Stat struct {
	ID             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Total_duration int64  `json:"total_duration"`
}

type ExtendedTaskWithStats struct {
	ExtendedTask  []ExtendedTask `json:"extended_task"`
	TaskStatistic []Stat         `json:"task_statistic"`
}

type TaskService interface {
	GetAll(userId int) (ExtendedTaskWithStats, error)
	Get(id int) (Task, error)
	Create(task taskDTO.Task) (Task, error)
	DeleteAllTasks() error
	CreateMany(task []taskDTO.Task) error
}
