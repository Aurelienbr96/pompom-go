//   Product Api:
//    version: 0.1
//    title: Product Api
//   Schemes: http, https
//   Host:
//   BasePath: /api/v1
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   SecurityDefinitions:
//    Bearer:
//     type: apiKey
//     name: Authorization
//     in: header
//   swagger:meta
package main

import (
	"fmt"
	"io"
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

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	// Set the content type to YAML to instruct the browser/client on how to interpret the file
	w.Header().Set("Content-Type", "text/plain")

	// Open the swagger.yaml file
	swaggerFile, err := os.Open("./swagger.yaml") // Make sure to provide the correct path
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer swaggerFile.Close()

	// Copy the contents of the file to the response writer
	_, err = io.Copy(w, swaggerFile)
	if err != nil {
		http.Error(w, "Error reading swagger file.", http.StatusInternalServerError)
	}
}

func setJsonApplicationHeader(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

func enableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

		taskService := services.NewTaskService(db)
		taskController := controllers.NewTaskController(taskService)

		mux.Handle("GET /task/{userId}", enableCors(setJsonApplicationHeader(taskController.GetAllTasksController)))
		mux.Handle("GET /task", enableCors(setJsonApplicationHeader(taskController.GetTask)))
		mux.Handle("POST /task", enableCors(setJsonApplicationHeader(taskController.CreateTask)))

		mux.Handle("GET /tag/{userId}", enableCors(setJsonApplicationHeader(TagController.GetAllTags)))
		mux.Handle("POST /tag", enableCors(setJsonApplicationHeader(TagController.CreateNewTag)))
		mux.Handle("/doc", enableCors(setJsonApplicationHeader(serveSwagger)))
		/* mux.HandleFunc("POST /task", controllers.PostTask) */

		if err := http.ListenAndServe("localhost:"+os.Getenv("PORT"), mux); err != nil {
			fmt.Println("error", err.Error())
		}
	}
}
