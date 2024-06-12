package model

type TagToTask struct {
	TaskId string `db:"taskid"`
	TagId  int    `db:"tagid"`
}
