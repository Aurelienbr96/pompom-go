package dto

import "time"

type Task struct {
	Name        string
	Description string
	Duration    int64
	TagId       int64
	Date        time.Time
}
