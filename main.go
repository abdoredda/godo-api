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

func main() {
	godotenv.Load()
	db, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database with err %v", err)
	}

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)

	taskHandler := handler.NewTaskHandler(taskService)
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/tasks", taskHandler.HandleTasks)
	http.HandleFunc("/tasks/{id}", taskHandler.HandleTask)
	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/users/{id}", userHandler.HandleUser)
	http.HandleFunc("/login", userHandler.HandleLogin)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to connect to the local server on port 8080 with err %v", err)
	}
}
