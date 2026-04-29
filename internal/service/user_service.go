package service

import (
	"errors"
	"hyweb-api/internal/middleware"
	"hyweb-api/internal/model"
	"hyweb-api/internal/repository"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (string, error)
	ChangePassword(email string, req model.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(req model.RegisterRequest) error {
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *userService) Login(req model.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)

	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	log.Printf("DB Password: [%s]", user.Password)
	log.Printf("Input Password: [%s]", req.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("Bcrypt Error: %v", err)
		return "", errors.New("invalid email or password")
	}

	token, err := middleware.GenerateToken(user.Email, user.Updated)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) ChangePassword(email string, req model.ChangePasswordRequest) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("old password incorrect")
	}

	if req.OldPassword == req.NewPassword {
		return errors.New("new password cannot be the same as old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(email, string(hashedPassword))
}
