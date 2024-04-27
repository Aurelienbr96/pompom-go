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

		mux.HandleFunc("GET /task", taskController.GetAllTasksController)
		mux.HandleFunc("GET /task/{id}", taskController.GetTask)
		mux.HandleFunc("POST /task", taskController.CreateTask)

		mux.HandleFunc("GET /tag", TagController.GetAllTags)
		mux.HandleFunc("/doc", serveSwagger)
		/* mux.HandleFunc("POST /task", controllers.PostTask) */

		if err := http.ListenAndServe("localhost:"+os.Getenv("PORT"), mux); err != nil {
			fmt.Println("error", err.Error())
		}
	}
}
