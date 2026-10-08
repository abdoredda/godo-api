package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/db"
	"task-manager-api/handler"
	"task-manager-api/middleware"
	"task-manager-api/repository"
	"task-manager-api/service"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	db, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database with err %v", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found")
	}

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	taskService := service.NewTaskService(taskRepo)
	authService := service.NewAuthService(jwtSecret)
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

	server := &http.Server{
		Addr: ":8080",
	}

	// Run HTTP server separately.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	// Wait for Ctrl+C or SIGTERM.
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-signalCtx.Done()

	log.Println("shutdown signal received")

	// Give existing requests up to 5 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")

}
