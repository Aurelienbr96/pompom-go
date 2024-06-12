package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pompom/go/src/dto"
	model "pompom/go/src/model"
	service "pompom/go/src/services"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type TagController struct {
	TagService service.TagService
}

func NewTagController(s service.TagService) *TagController {
	return &TagController{TagService: s}
}

func (tc TagController) GetAllTags(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tags, err := tc.TagService.GetAllTags(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Res400Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"404"`
	Message  string `json:"message" example:"fail to get data"`
}

func (tc TagController) CreateNewTag(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	var decodedTag = dto.Tag{}
	err := json.NewDecoder(r.Body).Decode(&decodedTag)
	if err != nil {
		http.Error(w, "could not uncode the json", http.StatusInternalServerError)
	}
	var tagToCreate = &model.TagToCreate{
		Name:   decodedTag.Name,
		Color:  decodedTag.Color,
		UserId: decodedTag.UserId,
	}
	err = validate.Struct(tagToCreate)
	if err != nil {
		// validation failed
		w.WriteHeader(http.StatusBadRequest)
		data := Res400Struct{
			Status:   "FAILED",
			HTTPCode: http.StatusBadRequest,
			Message:  err.Error(),
		}

		res, _ := json.Marshal(data)
		w.Write(res)
		return
	}
	if err := tc.TagService.CreateNewTag(*tagToCreate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode("success"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tc TagController) UpdateTag(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("hereeee")
	var decodedTag = model.TagToUpdate{}
	err := json.NewDecoder(r.Body).Decode(&decodedTag)
	if err != nil {
		http.Error(w, "could not uncode the json", http.StatusInternalServerError)
	}

	tag, err := tc.TagService.UpdateTag(decodedTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, err := json.Marshal(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(res)
}

func (tc TagController) DeleteTag(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tags, err := tc.TagService.DeleteTag(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
