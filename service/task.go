package service

import (
	"fmt"
	"task-manager-api/model"
)

type taskRepository interface {
	CreateTask(task model.Task, userID int) (model.Task, error)
	GetTasks(filter model.TaskFilter, userID int) ([]model.Task, error)
	UpdateTask(id int, data model.UpdateTask, userID int) (model.Task, error)
	DeleteTask(id int, userID int) error
	GetTask(id int, userID int) (model.Task, error)
}

type TaskService struct {
	repo taskRepository
}

func NewTaskService(repo taskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(data model.Task, userID int) (model.Task, error) {
	task, err := s.repo.CreateTask(data, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("creating task: %w", err)
	}
	return task, nil
}
func (s *TaskService) UpdateTask(id int, data model.UpdateTask, userID int) (model.Task, error) {
	task, err := s.repo.UpdateTask(id, data, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("updating task with id %d: %w", id, err)
	}
	return task, nil
}
func (s *TaskService) DeleteTask(id int, userID int) error {
	err := s.repo.DeleteTask(id, userID)
	if err != nil {
		return fmt.Errorf("deleting task with id %d: %w", id, err)
	}
	return nil
}
func (s *TaskService) GetTask(id int, userID int) (model.Task, error) {
	task, err := s.repo.GetTask(id, userID)
	if err != nil {
		return model.Task{}, fmt.Errorf("getting task: %w", err)
	}
	return task, nil
}
func (s *TaskService) GetTasks(filter model.TaskFilter, userID int) ([]model.Task, error) {
	tasks, err := s.repo.GetTasks(filter, userID)
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	return tasks, nil
}
