package repository

import (
	"database/sql"
	"task-manager-api/model"
)

// var tasks []model.Task

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) CreateTask(task model.Task) (model.Task, error) {

	row := r.db.QueryRow("INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id", task.Title, task.Done)
	err := row.Scan(&task.ID)
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) UpdateTask(id int, data model.UpdateTask) (model.Task, error) {
	var updated model.Task
	row := r.db.QueryRow("UPDATE tasks SET title = $1, done = $2 WHERE id = $3 RETURNING id, title, done", *data.Title, *data.Done, id)
	err := row.Scan(&updated.ID, &updated.Title, &updated.Done)
	if err == sql.ErrNoRows {
		return model.Task{}, model.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	return updated, nil
}

func (r *TaskRepository) DeleteTask(id int) error {
	result, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) GetTask(id int) (model.Task, error) {
	var t model.Task
	row := r.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = $1", id)
	err := row.Scan(&t.ID, &t.Title, &t.Done)
	if err == sql.ErrNoRows {
		return model.Task{}, model.ErrNotFound
	}

	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) GetTasks(filter model.TaskFilter) ([]model.Task, error) {
	var tasks []model.Task
	rows, err := r.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
