package model

import (
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
	UserId      int64     `db:"userid" json:"userid"`
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

type TaskToSend struct {
	Task
	Tag Tag `json:"tag"`
}

type ExtendedTaskWithStats struct {
	ExtendedTask  []ExtendedTask `json:"extended_task"`
	TaskStatistic []Stat         `json:"task_statistic"`
}
