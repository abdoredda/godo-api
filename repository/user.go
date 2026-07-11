package repository

import (
	"task-manager-api/model"
)

type UserRepository struct {
	db dbConnectionInterface
}

func NewUserRepository(db dbConnectionInterface) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
