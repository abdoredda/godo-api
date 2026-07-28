package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"task-manager-api/model"
)

type userServiceInterface interface {
	UserRegister(user model.User) (model.User, error)
	Login(user model.User) (model.User, string, error)
	GetUserById(id int) (model.User, error)
}
type UserHandler struct {
	userService userServiceInterface
}

func NewUserHandler(userService userServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (s *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		var userRequest model.UserRegisterRequest
		var user model.User
		var userRes model.UserResponse
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
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

		validationErrors := make(map[string]string)
		// username not empty
		if userRequest.Username == "" {
			validationErrors["username"] = "username field is required"
		}
		// password length <= 8
		if len(userRequest.Password) < 8 {
			validationErrors["password"] = "password must equal or greater than 8"
		}

		if userRequest.ConfirmPassword != userRequest.Password {
			validationErrors["confirm_password"] = "password and confirm_password fields must be the same"
		}

		if len(validationErrors) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]any{
				"errors": validationErrors,
			})

			if err != nil {
				log.Printf("error while trying to encode the response error = %v", err)
				return
			}
			return
		}

		user.Username = userRequest.Username
		user.Password = userRequest.Password

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

func (s *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		queryId := r.PathValue("id")
		if queryId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(queryId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error while trying to convert id to int error = %v", err)
			return
		}

		s.responseWithUser(id, w)
	}
}

func (s *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// parse and validate request (username & passwrd) $

	// compare request pass with db pass bcrypt
	// ture -> generate token and return user + token + err $
	// false -> return  401 (genaric error)
}

func (s *UserHandler) responseWithUser(id int, w http.ResponseWriter) {
	var user model.UserResponse
	data, err := s.userService.GetUserById(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.ID = data.ID
	user.Username = data.Username

	if err = json.NewEncoder(w).Encode(map[string]any{
		"message": "user fetched successfuly",
		"data":    user,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while trying to encode the response error = %v", err)
		return
	}
}
