package repository

import (
	"database/sql"
	"strconv"
	"task-manager-api/model"
)

// var tasks []model.Task

type TaskRepository struct {
	db dbConnectionInterface
}

type dbConnectionInterface interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Exec(string, ...any) (sql.Result, error)
}

func NewTaskRepository(db dbConnectionInterface) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) CreateTask(task model.Task, userID int) (model.Task, error) {

	row := r.db.QueryRow("INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id", task.Title, task.Done, userID)
	err := row.Scan(&task.ID)
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *TaskRepository) UpdateTask(id int, data model.UpdateTask, userID int) (model.Task, error) {
	var updated model.Task
	row := r.db.QueryRow("UPDATE tasks SET title = $1, done = $2 WHERE id = $3 AND user_id = $4 RETURNING id, title, done", *data.Title, *data.Done, id, userID)
	err := row.Scan(&updated.ID, &updated.Title, &updated.Done)
	if err == sql.ErrNoRows {
		return model.Task{}, model.ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	return updated, nil
}

func (r *TaskRepository) DeleteTask(id int, userID int) error {
	result, err := r.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
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

func (r *TaskRepository) GetTask(id int, userID int) (model.Task, error) {
	var t model.Task
	row := r.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = $1 AND user_id = $2 ", id, userID)
	err := row.Scan(&t.ID, &t.Title, &t.Done)
	if err == sql.ErrNoRows {
		return model.Task{}, model.ErrNotFound
	}

	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func (r *TaskRepository) GetTasks(filter model.TaskFilter, userID int) ([]model.Task, error) {
	var tasks []model.Task
	var args []any
	offset := (filter.Page - 1) * filter.Limit

	whereClause := "WHERE user_id = $1"
	args = append(args, userID)

	if filter.Done != nil {
		args = append(args, *filter.Done)
		whereClause += " AND done = $" + strconv.Itoa(len(args))
	}
	// offset =
	args = append(args, filter.Limit)
	limitPlaceholder := "$" + strconv.Itoa(len(args))

	args = append(args, offset)
	offsetPlaceholder := "$" + strconv.Itoa(len(args))

	query := "SELECT id, title, done FROM tasks " + whereClause + " ORDER BY id" +
		" LIMIT " + limitPlaceholder + " OFFSET " + offsetPlaceholder

	rows, err := r.db.Query(query, args...)
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
