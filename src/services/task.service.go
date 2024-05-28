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

func (c *TaskDb) GetAll(userId int) (taskModel.ExtendedTaskWithStats, error) {
	tasks := taskModel.ExtendedTaskWithStats{}
	err := c.DB.Select(&tasks.ExtendedTask, `
  SELECT 
    task.id, 
    task.name, 
    task.duration, 
    task.description, 
    tag.color, 
    task.tagId, 
    task.date 
  FROM 
    task 
  LEFT JOIN 
    tag 
  ON 
    tag.id = task.tagId
	WHERE 
		task.userid = $1
  ORDER BY 
    task.date DESC
`, userId)
	if err != nil {
		log.Fatalf("Error during query: %s", err)
	}

	var ids []int64
	for _, sec := range tasks.ExtendedTask {
		ids = append(ids, sec.TagId)
	}

	query, args, err := sqlx.In(`
		SELECT 
			tag.id as id, 
			tag.name, 
			SUM(task.duration) as Total_duration 
		FROM task 
		JOIN tag ON tag.id = task.tagId 
		WHERE tag.id IN (?) 
		GROUP BY tag.id, tag.name`, ids)
	if err != nil {
		log.Print("Error during query: %s", err)
	}
	query = c.DB.Rebind(query)
	err = c.DB.Select(&tasks.TaskStatistic, query, args...)
	if err != nil {
		log.Printf("Error during query: %s", err)
		return tasks, err
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
	sqlStatement := `INSERT INTO task (name, description, tagid, duration, date) VALUES (:name, :description, :tagid, :duration, :date) RETURNING *`

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
		Date:        task.Date,
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
