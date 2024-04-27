package services

import (
	"log"
	taskDTO "pompom/go/src/dto"
	taskModel "pompom/go/src/model"

	"github.com/jmoiron/sqlx"
)

type TaskDb struct {
	DB *sqlx.DB
}

func NewTaskService(db *sqlx.DB) taskModel.TaskService {
	return &TaskDb{DB: db}
}

func (c *TaskDb) GetAll() ([]taskModel.Task, error) {
	tasks := []taskModel.Task{}
	err := c.DB.Select(&tasks, "SELECT * FROM task")
	if err != nil {
		log.Fatalf("Error during query: %s", err)
	}
	return tasks, err
}

func (c *TaskDb) Get(id int) (taskModel.Task, error) {
	task := taskModel.Task{}
	err := c.DB.Get(&task, "SELECT * FROM task WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Error during query: %s", err)
	}
	// log.Printf("Task details: %+v", task)
	return task, err
}

func (c *TaskDb) Create(task taskDTO.Task) (taskModel.Task, error) {

	var newTask taskModel.Task
	sqlStatement := `INSERT INTO task (name, description, tag) VALUES (:name, :description, :tag) RETURNING *`

	namedStmt, err := c.DB.PrepareNamed(sqlStatement)
	if err != nil {
		log.Printf("Failed to prepare named statement: %s", err)
		return taskModel.Task{}, err
	}
	defer namedStmt.Close()

	err = namedStmt.Get(&newTask, taskDTO.Task{
		Name:        task.Name,
		Description: task.Description,
		Duration:    task.Duration,
		TagId:       task.TagId,
	})
	if err != nil {
		log.Printf("Failed to execute named statement: %s", err)
		return taskModel.Task{}, err
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
