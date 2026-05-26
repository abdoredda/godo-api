package service

import (
	"fmt"
	"task-manager-api/model"
)

type taskRepository interface {
	CreateTask(task model.Task) (model.Task, error)
	UpdateTask(id int, data model.UpdateTask) (model.Task, error)
	DeleteTask(int) error
	GetTask(int) (task model.Task, err error)
	GetTasks() (tasks []model.Task, err error)
}

type TaskService struct {
	repo taskRepository
}

func NewTaskService(repo taskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(data model.Task) (model.Task, error) {
	task, err := s.repo.CreateTask(data)
	if err != nil {
		return model.Task{}, fmt.Errorf("creating task: %w", err)
	}
	return task, nil
}
func (s *TaskService) UpdateTask(id int, data model.UpdateTask) (model.Task, error) {
	task, err := s.repo.UpdateTask(id, data)
	if err != nil {
		return model.Task{}, fmt.Errorf("updating task with id %d: %w", id, err)
	}
	return task, nil
}
func (s *TaskService) DeleteTask(id int) error {
	err := s.repo.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("deleting task with id %d: %w", id, err)
	}
	return nil
}
func (s *TaskService) GetTask(id int) (model.Task, error) {
	task, err := s.repo.GetTask(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("getting task: %w", err)
	}
	return task, nil
}
func (s *TaskService) GetTasks() ([]model.Task, error) {
	tasks, err := s.repo.GetTasks()
	if err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	return tasks, nil
}
