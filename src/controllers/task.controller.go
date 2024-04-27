package controllers

import (
	"encoding/json"
	"net/http"
	dto "pompom/go/src/dto"
	model "pompom/go/src/model"
	"strconv"
)

func setJsonApplicationHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

type TaskController struct {
	Service model.TaskService
}

func NewTaskController(s model.TaskService) *TaskController {
	return &TaskController{Service: s}
}

func (tc TaskController) GetAllTasksController(w http.ResponseWriter, r *http.Request) {

	setJsonApplicationHeader(w)
	tasks, err := tc.Service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tc TaskController) GetTask(w http.ResponseWriter, r *http.Request) {

	setJsonApplicationHeader(w)
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

func (tc TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	setJsonApplicationHeader(w)
	var t dto.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdTask, err := tc.Service.Create(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
