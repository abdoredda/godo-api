package handler

import (
	"encoding/json"
	"log"
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
		var userRes model.UserResponse

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid user syntax",
			})
			if err != nil {
				log.Printf("error while trying to encode the response error = %v", err)
				return
			}
			return
		}

		var validationErrors map
		// username not empty
		if user.Username == "" {
			validationErrors["errors"]["username"] = "username field is required"
		}
		// password length <= 8
		if len(user.Password) < 8 {
			validationErrors["errors"]["password"] = "password must equal or greater than 8"
		}

		data, err := s.userService.UserRegister(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error while trying to encode the response error = %v", err)
			err := json.NewEncoder(w).Encode(map[string]string{
				"error": "could not register user",
			})
			if err != nil {
				log.Printf("error while trying to encode the response error = %v", err)
				return
			}
			return
		}

		userRes.ID = data.ID
		userRes.Username = data.Username

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(map[string]any{
			"message": "user created successfully",
			"user":    userRes,
		})
		if err != nil {
			log.Printf("error while trying to encode the response error = %v", err)
			return
		}
		return
	default:
	}

}
