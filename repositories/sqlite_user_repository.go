package repositories

import (
	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/models"
	"github.com/jmoiron/sqlx"
)

type SQLiteUserRepository struct {
	DB *sqlx.DB
}

func (r *SQLiteUserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Select(&users, "SELECT id, name FROM users")
	return users, err
}

func (r *SQLiteUserRepository) Create(user models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
	return err
}
