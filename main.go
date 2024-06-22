package main

import (
	"context"
	"log"
	"os"

	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/handlers"
	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/repositories"
	"github.com/Micky90-sys/go-gin-postgres-redis-kafka/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	_ "modernc.org/sqlite"
)

func main() {
	// Database setup
	var userRepo repositories.UserRepositoryInterface

	if os.Getenv("DATABASE_URL") != "" {
		dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Printf("Unable to connect to PostgreSQL database: %v\n", err)
		} else {
			userRepo = &repositories.UserRepository{DB: dbpool}
		}
	}

	// Fallback to SQLite if PostgreSQL connection fails
	if userRepo == nil {
		log.Println("Using SQLite for local testing")
		sqliteDB, err := sqlx.Open("sqlite", ":memory:")
		if err != nil {
			log.Fatalf("Unable to connect to SQLite: %v\n", err)
		}
		defer sqliteDB.Close()

		_, err = sqliteDB.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
		if err != nil {
			log.Fatalf("Unable to create table in SQLite: %v\n", err)
		}

		userRepo = &repositories.SQLiteUserRepository{DB: sqliteDB}
	}

	// Redis setup
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	defer redisClient.Close()

	// Kafka setup
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_ADDR")},
		Topic:   "user_topic",
	})
	defer kafkaWriter.Close()

	// Service setup
	userService := &services.UserService{Repo: userRepo}

	// Handler setup
	userHandler := &handlers.UserHandler{Service: userService}

	// Router setup
	router := gin.Default()
	router.GET("/users", userHandler.GetUsers)
	router.POST("/users", userHandler.CreateUser)

	// Start server
	router.Run(":8080")
}
