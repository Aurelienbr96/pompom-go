package seed

import (
	"log"
	dto "pompom/go/src/dto"
	taskModel "pompom/go/src/model"
	services "pompom/go/src/services"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func generateTasks(row int) []dto.Task {
	createdTasks := make([]dto.Task, row)
	for i := 0; i < len(createdTasks); i++ {
		createdTasks[i] = dto.Task{
			Name: "name" + strconv.Itoa(i), Description: "description" + strconv.Itoa(i), TagId: int64(i + 1), Duration: int64(i),
		}
	}
	return createdTasks
}

func generateTag(name string, color string) dto.Tag {

	createdTag := dto.Tag{
		Name: name, Color: color,
	}
	return createdTag
}

type TaskService struct {
	Service taskModel.TaskService
	DB      *sqlx.DB
}

func NewSeedService(s taskModel.TaskService, db *sqlx.DB) *TaskService {
	return &TaskService{DB: db, Service: s}
}

func (c TaskService) CreateTasks(i int) {
	tasks := generateTasks(i)
	log.Printf("Task details: %+v", tasks)
	taskService := services.NewTaskService(c.DB)
	err := taskService.CreateMany(tasks)
	if err != nil {
		log.Printf("error:  %+v", err)
	}
	log.Printf("Tasks")
}

func (c TaskService) CreateTags() {
	reactNative := generateTag("React Native", "#B6C867")
	golang := generateTag("Golang", "#00CECB")
	vuejs := generateTag("vuejs", "#1B998B")
	algorithm := generateTag("algorithm", "#534D56")
	reactjs := generateTag("reactjs", "#FCB0B3")
	tags := [...]dto.Tag{reactNative, golang, vuejs, algorithm, reactjs}
	tagsSlice := tags[:]
	tagService := services.NewTagService(c.DB)
	err := tagService.CreateManyTags(tagsSlice)
	if err != nil {
		log.Printf("error:  %+v", err)
	}
	log.Printf("Tags created")
}

func (c TaskService) DeleteTasks() {
	err := c.Service.DeleteAllTasks()
	if err != nil {
		log.Printf("Failed to delete tasks, error: %v", err)
	}
	log.Printf("Deleted all tasks from db")
}

// TODO create func to create table task and tag with correct relations then update the create method
// create end points for tag
// create tests
// auth ?

/*

tag (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	color VARCHAR(255) NOT NULL
)


task (
	id SERIAL PRIMARY KEY,
	duration INT NOT NULL,
	name VARCHAR(255),
	description VARCHAR(255),
	tagId INT,
	FOREIGN KEY (tagId) REFERENCES tag(id)
)
*/

func (c TaskService) CreateDatabase() {
	_, err := c.DB.Exec(`CREATE TABLE tag (
							id SERIAL PRIMARY KEY,
							name VARCHAR(255) NOT NULL, 
							color VARCHAR(255) NOT NULL 
						  )`)
	if err != nil {
		log.Printf("An error happened when trying to create tag table  %+v", err)
	}
	_, err2 := c.DB.Exec(`CREATE TABLE task (
					id SERIAL PRIMARY KEY,
					duration INT NOT NULL, 
					name VARCHAR(255), 
					date DATE NOT NULL,
					description VARCHAR(255), 
					tagId INT, 
					FOREIGN KEY (tagId) REFERENCES tag(id))`)
	if err2 != nil {
		log.Printf("An error happened when trying to create task table  %+v", err)
	}
}

func (c TaskService) DeleteBd() {
	table := [2]string{"task", "tag"}
	for i := 0; i < len(table); i++ {
		_, err := c.DB.Exec("DROP TABLE " + table[i] + ";")
		if err != nil {
			log.Printf("failed to delete all database")
		}
	}
	log.Printf("database destroyed")
}
