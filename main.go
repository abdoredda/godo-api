package main

import (
	"log"
	"net/http"
	"task-manager-api/handler"
	"task-manager-api/repository"
	"task-manager-api/service"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {}

func main() {
	repo := repository.NewTaskRepository()
	taskService := service.NewTaskService(repo)
	h := handler.NewTaskHandler(taskService)

	http.HandleFunc("/tasks", h.HandleTasks)
	http.HandleFunc("/tasks/{id}", h.HandleTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to connect to the local server on port 8080 with err %v", err)
	}
}
