package model

import (
	taskDTO "pompom/go/src/dto"
)

/*
id INT PRIMARY KEY,
					duration INT NOT NULL,
					name VARCHAR(255),
					description VARCHAR(255),
					tagId INT,
					FOREIGN KEY (tagId) REFERENCES tag(id))
*/

type Task struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Duration    int64  `db:"duration" json:"duration"`
	Description string `db:"description" json:"description"`
	TagId       int64  `db:"tagid" json:"tagid"`
}

type TaskService interface {
	GetAll() ([]Task, error)
	Get(id int) (Task, error)
	Create(task taskDTO.Task) (Task, error)
	DeleteAllTasks() error
	CreateMany(task []taskDTO.Task) error
}
