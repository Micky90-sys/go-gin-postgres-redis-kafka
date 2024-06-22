package services

import (
	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/models"
	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/repositories"
)

type UserService struct {
	Repo repositories.UserRepositoryInterface
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) CreateUser(user models.User) error {
	return s.Repo.Create(user)
}
