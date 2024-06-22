package repositories

import (
	"context"

	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	Create(user models.User) error
}

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Create(user models.User) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO users (name) VALUES ($1)", user.Name)
	return err
}
