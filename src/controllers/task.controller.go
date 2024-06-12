package controllers

import (
	"encoding/json"
	"net/http"
	"pompom/go/src/dto"
	"pompom/go/src/model"
	service "pompom/go/src/services"
	"strconv"
	"time"
)

type FinalTask struct {
	model.Task
	Tags []model.Tag `json:"tags"`
}

type ErrorString struct {
	Error string `json:"error"`
}

func formatTasks(tasks []model.TaskToSend) interface{} {
	if len(tasks) == 0 {
		return ErrorString{
			Error: "You do not have any tasks",
		}
	}

	finaltasks := []FinalTask{}

	for _, val := range tasks {
		found := false
		for a := 0; a < len(finaltasks); a++ {
			if val.ID == finaltasks[a].Task.ID {
				finaltasks[a].Tags = append(finaltasks[a].Tags, val.Tag)
				found = true
				break
			}
		}
		if !found {
			newTask := FinalTask{
				Task: model.Task{
					ID: val.ID,
					TaskToCreate: model.TaskToCreate{
						Name:        val.Name,
						Date:        val.Date,
						Duration:    val.Duration,
						Description: val.Description,
					},
				},
				Tags: []model.Tag{val.Tag},
			}
			finaltasks = append(finaltasks, newTask)
		}
	}

	return finaltasks
}

type TaskController struct {
	TaskService      service.TaskService
	TagToTaskService service.TagToTaskService
}

func NewTaskController(taskService service.TaskService, tagToTaskService service.TagToTaskService) *TaskController {
	return &TaskController{TaskService: taskService, TagToTaskService: tagToTaskService}
}

func (tc TaskController) GetAllTasksController(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	tasks, err := tc.TaskService.GetAll(service.TaskParams{UserID: userId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	orderedTask := formatTasks(tasks)
	if err := json.NewEncoder(w).Encode(orderedTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tc TaskController) GetTask(w http.ResponseWriter, r *http.Request) {

	userIdStr := r.PathValue("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	taskIdStr := r.PathValue("taskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	task, err := tc.TaskService.GetAll(service.TaskParams{UserID: userId, TaskID: &taskId})
	orderedTask := formatTasks(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(orderedTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type TaskToCreate struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description,omitempty"`
	Duration    int64  `json:"duration" validate:"required,min=1"`
	TagId       []int  `json:"tagId"`
	Date        int64  `json:"date"`
	UserId      int    `json:"userid"`
}

func ConvertMilliToTime(milliseconds int64) time.Time {
	return time.Unix(milliseconds/1000, (milliseconds%1000)*1000000)
}

func (tc TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var t TaskToCreate
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskDate := ConvertMilliToTime(t.Date)
	taskToCreate := dto.Task{
		Date:        taskDate,
		Name:        t.Name,
		Description: t.Description,
		Duration:    t.Duration,
		UserId:      t.UserId,
	}

	createdTask, err := tc.TaskService.Create(taskToCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tagToTask, err := tc.TagToTaskService.CreateTagToTask(createdTask, t.TagId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(tagToTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
