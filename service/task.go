package service

import (
	"context"
	"fmt"
	"task-manager-api/model"
)

type taskRepository interface {
	CreateTask(ctx context.Context, task model.Task, userID int) (model.Task, error)
	GetTasks(ctx context.Context, filter model.TaskFilter, userID int) ([]model.Task, error)
	UpdateTask(ctx context.Context, id int, data model.UpdateTask, userID int) (model.Task, error)
	DeleteTask(ctx context.Context, id int, userID int) error
	GetTask(ctx context.Context, id int, userID int) (model.Task, error)
}

type TaskService struct {
	repo taskRepository
}

func NewTaskService(repo taskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, data model.Task, userID int) (model.Task, error) {
	task, err := s.repo.CreateTask(ctx, data, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("creating task: %w", err)
	}
	return task, nil
}
func (s *TaskService) UpdateTask(ctx context.Context, id int, data model.UpdateTask, userID int) (model.Task, error) {
	task, err := s.repo.UpdateTask(ctx, id, data, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("updating task with id %d: %w", id, err)
	}
	return task, nil
}
func (s *TaskService) DeleteTask(ctx context.Context, id int, userID int) error {
	err := s.repo.DeleteTask(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("deleting task with id %d: %w", id, err)
	}
	return nil
}
func (s *TaskService) GetTask(ctx context.Context, id int, userID int) (model.Task, error) {
	task, err := s.repo.GetTask(ctx, id, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("getting task: %w", err)
	}
	return task, nil
}
func (s *TaskService) GetTasks(ctx context.Context, filter model.TaskFilter, userID int) ([]model.Task, error) {
	tasks, err := s.repo.GetTasks(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	return tasks, nil
}
