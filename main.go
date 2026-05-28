package main

import (
	"log"
	"net/http"
	"os"
	"task-manager-api/db"
	"task-manager-api/handler"
	"task-manager-api/repository"
	"task-manager-api/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {}

func main() {
	godotenv.Load()
	db, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database with err %v", err)
	}
	repo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(repo)
	h := handler.NewTaskHandler(taskService)

	http.HandleFunc("/tasks", h.HandleTasks)
	http.HandleFunc("/tasks/{id}", h.HandleTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to connect to the local server on port 8080 with err %v", err)
	}
}
