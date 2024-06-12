package controllers

import (
	"encoding/json"
	"net/http"
	service "pompom/go/src/services"
	"strconv"
)

type StatisticController struct {
	StatisticService service.StatisticService
}

func NewStatisticController(s service.StatisticService) *StatisticController {
	return &StatisticController{StatisticService: s}
}

func (tc StatisticController) GetAllStatistics(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("userId")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tags, err := tc.StatisticService.GetStatistic(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
