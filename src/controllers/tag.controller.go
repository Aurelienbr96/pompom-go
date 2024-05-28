package controllers

import (
	"encoding/json"
	"net/http"
	"pompom/go/src/dto"
	model "pompom/go/src/model"
	"strconv"
)

type TagController struct {
	Service model.TagService
}

func NewTagController(s model.TagService) *TagController {
	return &TagController{Service: s}
}

// swagger:operation GET /tag Tag getAllTags
// Get Tag List
//
// ---
// responses:
//
//  401: CommonError
//  200: CommonSuccess
func (tc TagController) GetAllTags(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tags, err := tc.Service.GetAllTags(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (tc TagController) CreateNewTag(w http.ResponseWriter, r *http.Request) {
	var newTag = dto.Tag{}
	err := json.NewDecoder(r.Body).Decode(&newTag)
	if err != nil {
		http.Error(w, "could not uncode the json", http.StatusInternalServerError)
	}
	if err := tc.Service.CreateNewTag(newTag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode("success"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
