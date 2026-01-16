package services

import (
	"sales-management-api/internal/models"
	"sales-management-api/internal/repositories"
)

type UserService struct {
	users *repositories.UserRepo
}

func NewUserService(users *repositories.UserRepo) *UserService {
	return &UserService{users: users}
}

func (s *UserService) List() ([]models.User, error) {
	return s.users.FindAll()
}

func (s *UserService) Detail(id uint) (*models.User, error) {
	return s.users.FindByID(id)
}
