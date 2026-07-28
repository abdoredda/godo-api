package service

import (
	"errors"
	"fmt"
	"log"
	"task-manager-api/model"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	CreateUser(user model.User) (model.User, error)
	GetUserByUsername(username string) (model.User, error)
	GetUserById(id int) (model.User, error)
}

type UserService struct {
	UserRepo    userRepository
	AuthService *AuthService
}

func NewUserService(userRepository userRepository, authService *AuthService) *UserService {
	return &UserService{UserRepo: userRepository, AuthService: authService}
}

func (s *UserService) UserRegister(user model.User) (model.User, error) {

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("error when trying to hash the password error = %v", err)
	}

	u := model.User{Username: user.Username, Password: string(hashedPassword)}
	data, err := s.UserRepo.CreateUser(u)
	if err != nil {
		return model.User{}, fmt.Errorf("register service: error when trying to create new user error = %v", err)
	}

	return data, nil
}

func (s *UserService) Login(user model.User) (model.User, string, error) {

	data, err := s.UserRepo.GetUserByUsername(user.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Printf("login attempt for nonexistent user username = %v", user.Username)
			return model.User{}, "", model.ErrInvalidCredentials
		}
		return model.User{}, "", fmt.Errorf("login service: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password))
	if err != nil {
		log.Printf("failed login attempt for username %q: bad password", user.Username)
		return model.User{}, "", model.ErrInvalidCredentials
	}

	token, err := s.AuthService.GenerateToken(data.ID)
	if err != nil {
		return model.User{}, "", fmt.Errorf("login service: %w", err)
	}

	return data, token, nil
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	data, err := s.UserRepo.GetUserById(id)
	if err != nil {
		return model.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return data, nil
}
