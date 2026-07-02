package repository

import (
	"database/sql"
	"reflect"
	"task-manager-api/model"
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

func TestGetTasks_FilterByDone(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("Error when trying to start db transaction error = %v", err)
	}
	defer tx.Rollback()

	task_1 := model.Task{
		Title: "Learn Go",
		Done:  true,
	}
	task_2 := model.Task{
		Title: "T2",
		Done:  true,
	}
	task_3 := model.Task{
		Title: "T3",
		Done:  false,
	}

	query := "INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id;"
	err = tx.QueryRow(query, task_1.Title, task_1.Done).Scan(&task_1.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}
	err = tx.QueryRow(query, task_2.Title, task_2.Done).Scan(&task_2.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}
	err = tx.QueryRow(query, task_3.Title, task_3.Done).Scan(&task_3.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}

	repo := NewTaskRepository(tx)
	done := true
	undone := false

	tests := []struct {
		name           string
		filterInp      model.TaskFilter
		expectedOutput []model.Task
	}{
		{
			name:      "done not provided",
			filterInp: model.TaskFilter{Page: 1, Limit: 10},
			expectedOutput: []model.Task{
				task_1,
				task_2,
				task_3,
			},
		},
		{
			name:      "done = true",
			filterInp: model.TaskFilter{Page: 1, Limit: 10, Done: &done},
			expectedOutput: []model.Task{
				task_1,
				task_2,
			},
		},
		{
			name:      "done = false",
			filterInp: model.TaskFilter{Page: 1, Limit: 10, Done: &undone},
			expectedOutput: []model.Task{
				task_3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := repo.GetTasks(test.filterInp)
			if err != nil {
				t.Fatalf("error when calling GetTasks() error = %v", err)
			}
			if !reflect.DeepEqual(data, test.expectedOutput) {
				t.Errorf("return %v expected %v", data, test.expectedOutput)
			}
		})
	}

}
