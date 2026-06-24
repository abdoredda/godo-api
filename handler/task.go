package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"task-manager-api/model"
)

type taskService interface {
	CreateTask(task model.Task) (model.Task, error)
	UpdateTask(id int, data model.UpdateTask) (model.Task, error)
	DeleteTask(int) error
	GetTask(int) (task model.Task, err error)
	GetTasks(model.TaskFilter) (tasks []model.Task, err error)
}

type TaskHandler struct {
	taskService taskService
}

func NewTaskHandler(service taskService) *TaskHandler {
	return &TaskHandler{taskService: service}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page := 1
	var done *bool
	limit := 10

	switch r.Method {
	case "GET":
		if s := r.URL.Query().Get("page"); s != "" {
			value, err := strconv.Atoi(s)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "page must be a valid number",
				}); err != nil {
					log.Printf("failed to encode error response: %v", err)
				}
				return
			}
			if value <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "page must be greater than 0",
				}); err != nil {
					log.Printf("failed to encode error response: %v", err)
				}
				return
			}
			page = value
		}

		if s := r.URL.Query().Get("limit"); s != "" {
			value, err := strconv.Atoi(s)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "limit must be a valid number",
				}); err != nil {
					log.Printf("failed to encode error response: %v", err)
				}
				return
			}
			if value <= 0 || value > 100 {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "limit must be greater than 0 and less than or equal 100",
				}); err != nil {
					log.Printf("failed to encode error response: %v", err)
				}
				return
			}
			limit = value
		}

		if d := r.URL.Query().Get("done"); d != "" {
			value, err := strconv.ParseBool(d)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "done must be a valid boolean value",
				}); err != nil {
					log.Printf("failed to encode error response: %v", err)
				}
				return
			}

			done = &value
		}

		filter := model.TaskFilter{
			Done:  done,
			Page:  page,
			Limit: limit,
		}

		tasks, err := h.taskService.GetTasks(filter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "error while fetching tasks",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			log.Printf("failed to encode tasks: %v", err)
		}

	case "POST":
		var rTask model.Task
		if err := json.NewDecoder(r.Body).Decode(&rTask); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid task syntax",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		var validationErr []string
		if rTask.Title == "" {
			validationErr = append(validationErr, "title must not be empty")
		} else if len(rTask.Title) < 3 {
			validationErr = append(validationErr, "title must be at least 3 characters")
		}
		if len(validationErr) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string][]string{
				"error": validationErr,
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		task, err := h.taskService.CreateTask(rTask)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "error while creating the task",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(task); err != nil {
			log.Printf("failed to encode task response: %v", err)
		}
		// ---
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		}); err != nil {
			log.Printf("failed to encode error response: %v", err)
		}
		return
	}
}

func (h *TaskHandler) HandleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	trimmedTaskID := r.PathValue("id")
	taskId, err := strconv.Atoi(trimmedTaskID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "can't parsing id string",
		}); err != nil {
			log.Printf("failed to encode error response: %v", err)
		}
		return
	}

	switch r.Method {
	case "GET":
		task, err := h.taskService.GetTask(taskId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "error while getting the task",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		json.NewEncoder(w).Encode(task)
		return

	case "PUT":
		var data model.UpdateTask
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "can't parsing the request body",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		var validationErr []string
		if data.Title == nil || *data.Title == "" {
			validationErr = append(validationErr, "title must not be empty")
		} else if len(*data.Title) < 3 {
			validationErr = append(validationErr, "title must be at least 3 characters")
		}
		if data.Done == nil {
			validationErr = append(validationErr, "done must be provided")
		}
		if len(validationErr) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string][]string{
				"error": validationErr,
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		task, err := h.taskService.UpdateTask(taskId, data)

		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "error while updating the task",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		json.NewEncoder(w).Encode(task)
		return

	case "DELETE":

		err = h.taskService.DeleteTask(taskId)

		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "error while deleting the task",
			}); err != nil {
				log.Printf("failed to encode error response: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		}); err != nil {
			log.Printf("failed to encode error response: %v", err)
		}
		return
	}
}
