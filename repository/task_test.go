package repository

import (
	"database/sql"
	"testing"
)

func TestGetTask_Found(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction error: %v", err)
	}
	defer tx.Rollback()
	taskTitle := "testTask"
	taskDone := false
	var id int
	query := "INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id;"

	err = tx.QueryRow(query, taskTitle, taskDone).Scan(&id)
	if err == sql.ErrNoRows {
		t.Fatalf("No rows found returned error: %v", err)
	}

	repo := NewTaskRepository(tx)
	task, err := repo.GetTask(id)
	if err != nil {
		t.Fatalf("GetTask(%d) returned error: %v", id, err)
	}

	if task.Title != taskTitle || task.Done != taskDone {
		t.Errorf("title = %v, done = %v wanted title = %v, done = %v", task.Title, task.Done, taskTitle, taskDone)
	}

}
