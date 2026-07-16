package model

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type UpdateTask struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

type TaskFilter struct {
	Page  int
	Limit int
	Done  *bool
}
