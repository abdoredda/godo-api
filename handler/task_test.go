package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"task-manager-api/middleware"
	"task-manager-api/model"
	"task-manager-api/service"
	"testing"
)

type taskRepo struct {
	tasks []model.Task
	err   error
}

func (s *taskRepo) GetTasks(ctx context.Context, filter model.TaskFilter, userID int) ([]model.Task, error) {
	return s.tasks, s.err
}

func (s *taskRepo) CreateTask(ctx context.Context, task model.Task, userID int) (model.Task, error) {
	return model.Task{}, nil
}
func (s *taskRepo) UpdateTask(ctx context.Context, id int, data model.UpdateTask, userID int) (model.Task, error) {
	return model.Task{}, nil
}

func (s *taskRepo) DeleteTask(ctx context.Context, id int, userID int) error {
	return nil
}

func (s *taskRepo) GetTask(ctx context.Context, id int, userID int) (model.Task, error) {
	return model.Task{}, nil
}

func TestGetTasks_RepoErr(t *testing.T) {

	taskService := service.NewTaskService(&taskRepo{err: errors.New("error")})
	taskHandler := NewTaskHandler(taskService)
	authService := service.NewAuthService("test-secret")
	tasksHandler := http.HandlerFunc(taskHandler.HandleTasks)

	token, err := authService.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	protectedTasksHandler := middleware.AuthMiddleware(tasksHandler, authService)

	// adapt taskHandler to satisfy handler interface to can pass it to handl() for registeration
	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		strings.NewReader(``),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	protectedTasksHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestGetTasks_Success(t *testing.T) {

	expectedTasks := []model.Task{
		{ID: 1, Title: "First task", Done: false},
		{ID: 2, Title: "Second task", Done: true},
	}
	taskService := service.NewTaskService(&taskRepo{tasks: expectedTasks})
	taskHandler := NewTaskHandler(taskService)
	authService := service.NewAuthService("test-secret")
	tasksHandler := http.HandlerFunc(taskHandler.HandleTasks)

	token, err := authService.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	protectedTasksHandler := middleware.AuthMiddleware(tasksHandler, authService)
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	protectedTasksHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var gotTasks []model.Task
	if err := json.NewDecoder(rec.Body).Decode(&gotTasks); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if !reflect.DeepEqual(gotTasks, expectedTasks) {
		t.Errorf("response tasks = %#v, want %#v", gotTasks, expectedTasks)
	}

}
