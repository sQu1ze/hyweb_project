package repository

import (
	"database/sql"
	"errors"
	"hyweb-api/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	UpdatePassword(email string, newPassword string) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(user *model.User) error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	_, err := r.db.Exec(query, user.Email, user.Password)
	return err
}

func (r *mysqlUserRepository) FindByEmail(email string) (*model.User, error) {
	query := "SELECT email, password, created, updated FROM users WHERE email = ?"
	row := r.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.Email, &user.Password, &user.Created, &user.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) UpdatePassword(email string, newPassword string) error {
	query := "UPDATE users SET password = ?, updated = CURRENT_TIMESTAMP WHERE email = ?"
	_, err := r.db.Exec(query, newPassword, email)
	return err
}
