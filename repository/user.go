package repository

import (
	"database/sql"
	"task-manager-api/model"
)

type UserRepository struct {
	db dbConnectionInterface
}

func NewUserRepository(db dbConnectionInterface) *UserRepository {
	return &UserRepository{db: db}
}

func (s *UserRepository) CreateUser(user model.User) (model.User, error) {
	err := s.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserRepository) GetUserById(id int) (model.User, error) {
	var user model.User
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, model.ErrNotFound
		}
		return model.User{}, err
	}
	return user, nil
}
