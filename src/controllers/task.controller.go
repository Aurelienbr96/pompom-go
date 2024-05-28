package controllers

import (
	"encoding/json"
	"net/http"
	"pompom/go/src/dto"
	model "pompom/go/src/model"
	"strconv"
	"time"
)

type TaskController struct {
	Service model.TaskService
}

func NewTaskController(s model.TaskService) *TaskController {
	return &TaskController{Service: s}
}

func (tc TaskController) GetAllTasksController(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	tasks, err := tc.Service.GetAll(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tc TaskController) GetTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	task, err := tc.Service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type TaskToCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Duration    int64  `json:"duration"`
	TagId       int64  `json:"tagId"`
	Date        int64  `json:"date"`
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
		TagId:       t.TagId,
	}
	createdTask, err := tc.Service.Create(taskToCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
