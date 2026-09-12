package repository

import (
	"errors"
	"reflect"
	"task-manager-api/model"
	"testing"
)

func createTestUser(t *testing.T, db dbConnectionInterface, username string) int {
	t.Helper()
	var userID int
	err := db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, "password").Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return userID
}

func TestGetTask_Found(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction error: %v", err)
	}
	defer tx.Rollback()

	userID := createTestUser(t, tx, "user_get_task")
	taskTitle := "testTask"
	taskDone := false
	var id int
	query := "INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id;"

	err = tx.QueryRow(query, taskTitle, taskDone, userID).Scan(&id)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	repo := NewTaskRepository(tx)
	task, err := repo.GetTask(id, userID)
	if err != nil {
		t.Fatalf("GetTask(%d, %d) returned error: %v", id, userID, err)
	}

	if task.Title != taskTitle || task.Done != taskDone {
		t.Errorf("title = %v, done = %v wanted title = %v, done = %v", task.Title, task.Done, taskTitle, taskDone)
	}
}

func TestGetTask_OtherUserNotFound(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction: %v", err)
	}
	defer tx.Rollback()

	ownerID := createTestUser(t, tx, "owner_user_get")
	otherUserID := createTestUser(t, tx, "other_user_get")

	var taskID int
	err = tx.QueryRow("INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id", "private task", false, ownerID).Scan(&taskID)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	repo := NewTaskRepository(tx)
	_, err = repo.GetTask(taskID, otherUserID)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound when getting another user's task, got: %v", err)
	}
}

func TestGetTasks_FilterByDone(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("Error when trying to start db transaction error = %v", err)
	}
	defer tx.Rollback()

	userID := createTestUser(t, tx, "user_filter")
	otherUserID := createTestUser(t, tx, "user_filter_other")

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

	query := "INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id;"
	err = tx.QueryRow(query, task_1.Title, task_1.Done, userID).Scan(&task_1.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}
	err = tx.QueryRow(query, task_2.Title, task_2.Done, userID).Scan(&task_2.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}
	err = tx.QueryRow(query, task_3.Title, task_3.Done, userID).Scan(&task_3.ID)
	if err != nil {
		t.Fatalf("Error when trying to insert a row into tasks table error %v = ", err)
	}

	// Insert task for another user to verify isolation
	var otherTaskID int
	err = tx.QueryRow(query, "Other user task", true, otherUserID).Scan(&otherTaskID)
	if err != nil {
		t.Fatalf("Error inserting other user task: %v", err)
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
			data, err := repo.GetTasks(test.filterInp, userID)
			if err != nil {
				t.Fatalf("error when calling GetTasks() error = %v", err)
			}
			if !reflect.DeepEqual(data, test.expectedOutput) {
				t.Errorf("return %v expected %v", data, test.expectedOutput)
			}
		})
	}
}

func TestDeleteTask_Success(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction: %v", err)
	}
	defer tx.Rollback()

	userID := createTestUser(t, tx, "user_del_success")
	var taskID int
	err = tx.QueryRow("INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id", "task to delete", false, userID).Scan(&taskID)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	repo := NewTaskRepository(tx)
	err = repo.DeleteTask(taskID, userID)
	if err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	_, err = repo.GetTask(taskID, userID)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got: %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction: %v", err)
	}
	defer tx.Rollback()

	userID := createTestUser(t, tx, "user_del_not_found")
	repo := NewTaskRepository(tx)
	err = repo.DeleteTask(999999, userID)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound for non-existent task, got: %v", err)
	}
}

func TestDeleteTask_OtherUserNotFound(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction: %v", err)
	}
	defer tx.Rollback()

	ownerID := createTestUser(t, tx, "owner_del")
	otherUserID := createTestUser(t, tx, "other_del")

	var taskID int
	err = tx.QueryRow("INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id", "owner task", false, ownerID).Scan(&taskID)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	repo := NewTaskRepository(tx)
	// Other user attempts to delete owner's task
	err = repo.DeleteTask(taskID, otherUserID)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound when deleting another user's task, got: %v", err)
	}

	// Verify task still exists for owner
	task, err := repo.GetTask(taskID, ownerID)
	if err != nil {
		t.Fatalf("expected owner task to still exist, got error: %v", err)
	}
	if task.ID != taskID {
		t.Errorf("expected task ID %d, got %d", taskID, task.ID)
	}
}

func TestUpdateTask_OtherUserNotFound(t *testing.T) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("error when trying to start db transaction: %v", err)
	}
	defer tx.Rollback()

	ownerID := createTestUser(t, tx, "owner_update")
	otherUserID := createTestUser(t, tx, "other_update")

	var taskID int
	err = tx.QueryRow("INSERT INTO tasks (title, done, user_id) VALUES ($1, $2, $3) RETURNING id", "original title", false, ownerID).Scan(&taskID)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	repo := NewTaskRepository(tx)
	newTitle := "hacked title"
	isDone := true
	_, err = repo.UpdateTask(taskID, model.UpdateTask{Title: &newTitle, Done: &isDone}, otherUserID)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("expected ErrNotFound when updating another user's task, got: %v", err)
	}

	// Verify owner task unchanged
	task, err := repo.GetTask(taskID, ownerID)
	if err != nil {
		t.Fatalf("expected owner task to exist, got error: %v", err)
	}
	if task.Title != "original title" || task.Done != false {
		t.Errorf("expected original task to be unmodified, got title=%s done=%v", task.Title, task.Done)
	}
}
