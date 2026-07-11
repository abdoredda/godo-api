package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/model"
)

type userServiceInterface interface {
	UserRegister(user model.User) (model.User, error)
}
type UserHandler struct {
	userService userServiceInterface
}

func (s *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		
	case "POST":
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
	}

}
