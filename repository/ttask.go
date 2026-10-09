package repository

import (
	"database/sql"
	"task-manager-api/model"
)

type TtaskRepo struct {
	db tdbConnectionInterface
}

type tdbConnectionInterface interface {
	QueryRow(string ...any) (*sql.Rows, error)
}

func NewTtaskRepo(db tdbConnectionInterface) *TtaskRepo {
	return &TtaskRepo{db}
}

func (r *TtaskRepo) GetTasks() []model.Task {
	return []model.Task{
		{ID: 1, Title: "T1", Done: false},
		{ID: 2, Title: "T2", Done: true},
	}
}
