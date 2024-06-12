package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	controllers "pompom/go/src/controllers"
	db "pompom/go/src/db"
	seed "pompom/go/src/seed"
	services "pompom/go/src/services"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setJsonApplicationHeader(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func enableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to fetch env variable: %v", err)
	}
	db, err := db.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close()

	if len(os.Args) > 1 {
		seed.Seed(db)

	} else {
		mux := http.NewServeMux()

		tagService := services.NewTagService(db)
		TagController := controllers.NewTagController(tagService)

		tagToTaskService := services.NewTaskToTagService(db)
		taskService := services.NewTaskService(db)
		taskController := controllers.NewTaskController(taskService, tagToTaskService)

		statisticService := services.NewStatistic(db)
		statisticController := controllers.NewStatisticController(statisticService)

		mux.Handle("PUT /tag", enableCors(setJsonApplicationHeader(TagController.UpdateTag)))

		mux.Handle("GET /task/{userId}", enableCors(setJsonApplicationHeader(taskController.GetAllTasksController)))
		mux.Handle("GET /task/{userId}/{taskId}", enableCors(setJsonApplicationHeader(taskController.GetTask)))
		mux.Handle("POST /task", enableCors(setJsonApplicationHeader(taskController.CreateTask)))

		mux.Handle("GET /statistic/{userId}", enableCors(setJsonApplicationHeader(statisticController.GetAllStatistics)))

		mux.Handle("GET /tag/{userId}", enableCors(setJsonApplicationHeader(TagController.GetAllTags)))
		mux.Handle("POST /tag", enableCors(setJsonApplicationHeader(TagController.CreateNewTag)))
		mux.Handle("DELETE /tag/{userId}", enableCors(setJsonApplicationHeader(TagController.DeleteTag)))

		mux.Handle("OPTIONS /tag", enableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		})))

		if err := http.ListenAndServe("localhost:"+os.Getenv("PORT"), mux); err != nil {
			fmt.Println("error", err.Error())
		}
	}
}
