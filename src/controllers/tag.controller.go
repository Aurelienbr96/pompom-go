package controllers

import (
	"encoding/json"
	"net/http"
	model "pompom/go/src/model"
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

	setJsonApplicationHeader(w)
	tags, err := tc.Service.GetAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
