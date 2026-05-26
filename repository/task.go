package repository

import (
	"task-manager-api/model"
)

var tasks []model.Task

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) CreateTask(task model.Task) (model.Task, error) {
	task.ID = len(tasks) + 1
	tasks = append(tasks, task)
	return task, nil
}

func (r *TaskRepository) UpdateTask(id int, data model.UpdateTask) (model.Task, error) {
	for i, task := range tasks {
		if id == task.ID {
			if data.Title != nil {
				tasks[i].Title = *data.Title
			}
			if data.Done != nil {
				tasks[i].Done = *data.Done
			}
			return tasks[i], nil
		}
	}
	return model.Task{}, model.ErrNotFound
}

func (r *TaskRepository) DeleteTask(id int) error {
	for i, task := range tasks {
		if id == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return model.ErrNotFound
}

func (r *TaskRepository) GetTask(id int) (model.Task, error) {
	for _, task := range tasks {
		if id == task.ID {
			return task, nil
		}
	}
	return model.Task{}, model.ErrNotFound
}

func (r *TaskRepository) GetTasks() ([]model.Task, error) {
	return tasks, nil
}
