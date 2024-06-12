package services

import (
	"log"
	taskDTO "pompom/go/src/dto"
	model "pompom/go/src/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskService interface {
	GetAll(params TaskParams) ([]model.TaskToSend, error)
	Get(id int) (model.Task, error)
	Create(task taskDTO.Task) (model.Task, error)
	DeleteAllTasks() error
	CreateMany(task []taskDTO.Task) error
}

type TaskDb struct {
	DB *sqlx.DB
}

func NewTaskService(db *sqlx.DB) TaskService {
	return &TaskDb{DB: db}
}

type TaskParams struct {
	UserID int
	TaskID *int
}

func (c *TaskDb) GetAll(params TaskParams) ([]model.TaskToSend, error) {
	queryBase := `
    SELECT 
        task.id AS task_id, 
        task.name AS task_name, 
        task.duration AS task_duration, 
        task.description AS task_description, 
				task.userid AS user_id,
        task.date AS task_date,
        tag.id AS tag_id,
        tag.name AS tag_name,
        tag.color AS tag_color,
        tag.userid AS tag_userid
    FROM 
        task 
    LEFT JOIN 
        tagtotask 
    ON 
        task.id = tagtotask.taskid
    LEFT JOIN
        tag
    ON
        tagtotask.tagid = tag.id
    WHERE 
        task.userid = $1
    `

	var queryParams []interface{}
	queryParams = append(queryParams, params.UserID)

	if params.TaskID != nil {
		queryBase += " AND task.id = $2"
		queryParams = append(queryParams, *params.TaskID)
	}

	queryBase += " ORDER BY task.date DESC"

	rows, err := c.DB.Queryx(queryBase, queryParams...)
	if err != nil {
		return []model.TaskToSend{}, err
	}
	defer rows.Close()

	var tasks []model.TaskToSend

	for rows.Next() {
		var taskTag struct {
			TaskID          int64     `db:"task_id"`
			TaskName        string    `db:"task_name"`
			TaskDuration    int64     `db:"task_duration"`
			TaskDescription string    `db:"task_description"`
			UserID          int64     `db:"user_id"`
			TaskDate        time.Time `db:"task_date"`
			TagID           int       `db:"tag_id"`
			TagName         string    `db:"tag_name"`
			TagColor        string    `db:"tag_color"`
			TagUserID       int64     `db:"tag_userid"`
		}

		err := rows.StructScan(&taskTag)
		if err != nil {
			return []model.TaskToSend{}, err
		}

		task := model.TaskToSend{
			Task: model.Task{
				ID: taskTag.TaskID,

				TaskToCreate: model.TaskToCreate{
					Name:        taskTag.TaskName,
					Date:        taskTag.TaskDate,
					Duration:    taskTag.TaskDuration,
					Description: taskTag.TaskDescription,
					UserId:      taskTag.UserID,
				},
			},
			Tag: model.Tag{
				ID: int64(taskTag.TagID),
				TagToCreate: model.TagToCreate{
					Name:  taskTag.TagName,
					Color: taskTag.TagColor,
				},
			},
		}
		tasks = append(tasks, task)

	}

	return tasks, nil
}

func (c *TaskDb) Get(id int) (model.Task, error) {
	task := model.Task{}
	err := c.DB.Get(&task, "SELECT * FROM task WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Error during query: %s", err)
	}
	// log.Printf("Task details: %+v", task)
	return task, err
}

func (c *TaskDb) Create(task taskDTO.Task) (model.Task, error) {

	var newTask model.Task
	sqlStatement := `INSERT INTO task (name, description, duration, date, userid) VALUES (:name, :description, :duration, :date, :userid) RETURNING *`

	namedStmt, err := c.DB.PrepareNamed(sqlStatement)
	if err != nil {
		log.Printf("Failed to prepare named statement: %s", err)
		return model.Task{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&newTask, taskDTO.Task{
		Name:        task.Name,
		Description: task.Description,
		Duration:    task.Duration,
		Date:        task.Date,
		UserId:      task.UserId,
	})
	if err != nil {
		log.Printf("Failed to execute named statement: %s", err)
		return model.Task{}, err
	}

	return newTask, nil
}

func (c *TaskDb) DeleteAllTasks() error {
	sqlStatement := "DELETE FROM task"
	_, err := c.DB.Exec(sqlStatement)
	if err != nil {
		log.Printf("Failed to prepare named statement: %s", err)
		return err
	}

	return err
}

func (c *TaskDb) CreateMany(tasks []taskDTO.Task) error {
	tx, err := c.DB.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed("INSERT INTO task (name, description, duration, tagid) VALUES (:name, :description, :duration, :tagid)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, task := range tasks {
		createdTask, err := stmt.Exec(task)
		if err != nil {
			tx.Rollback()
			return err
		}
		rowsAffected, err := createdTask.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("New task created, rows affected: %d", rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
