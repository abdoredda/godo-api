package service

import (
	"fmt"
	"task-manager-api/model"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	CreateUser(user model.User) (model.User, error)
}

type UserService struct {
	UserRepo userRepository
}

func NewUserService(userRepository userRepository) *UserService {
	return &UserService{UserRepo: userRepository}
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
		return model.User{}, fmt.Errorf("error when trying to create a user error = %v", err)
	}

	return data, nil
}
