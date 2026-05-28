package model

import "errors"

var ErrNotFound = errors.New("not found")

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type UpdateTask struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}
