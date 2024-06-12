package services

import (
	model "pompom/go/src/model"

	"github.com/jmoiron/sqlx"
)

type tagToTaskDb struct {
	DB *sqlx.DB
}

type TagToTaskService interface {
	CreateTagToTask(task model.Task, tagIds []int) (model.TagToTask, error)
}

func NewTaskToTagService(db *sqlx.DB) TagToTaskService {
	return &tagToTaskDb{DB: db}
}

func (c *tagToTaskDb) CreateTagToTask(task model.Task, tagIds []int) (model.TagToTask, error) {
	tx, err := c.DB.Beginx()
	if err != nil {
		return model.TagToTask{}, err
	}

	stmt, err := tx.PrepareNamed("INSERT INTO tagtotask (taskid, tagid) VALUES (:taskid, :tagid)")
	if err != nil {
		tx.Rollback()
		return model.TagToTask{}, err
	}
	defer stmt.Close()

	for _, tagId := range tagIds {
		params := map[string]interface{}{
			"taskid": task.ID,
			"tagid":  tagId,
		}
		_, err := stmt.Exec(params)
		if err != nil {
			tx.Rollback()
			return model.TagToTask{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.TagToTask{}, err
	}
	return model.TagToTask{}, nil
}
