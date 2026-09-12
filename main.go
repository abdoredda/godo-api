package main

import (
	"log"
	"net/http"
	"os"
	"task-manager-api/db"
	"task-manager-api/handler"
	"task-manager-api/middleware"
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
	authService := service.NewAuthService(os.Getenv("JWT_SECRET"))
	userService := service.NewUserService(userRepo, authService)

	taskHandler := handler.NewTaskHandler(taskService)
	tasksHandler := http.HandlerFunc(taskHandler.HandleTasks)
	taskItemHandler := http.HandlerFunc(taskHandler.HandleTask)
	protectedTasksHandler := middleware.AuthMiddleware(tasksHandler, authService)
	protectedTaskItemHandler := middleware.AuthMiddleware(taskItemHandler, authService)
	userHandler := handler.NewUserHandler(userService)

	http.Handle("/tasks", protectedTasksHandler)
	http.Handle("/tasks/{id}", protectedTaskItemHandler)
	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/users/{id}", userHandler.HandleUser)
	http.HandleFunc("/login", userHandler.HandleLogin)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to connect to the local server on port 8080 with err %v", err)
	}
}
