package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Ijne/homepage_app/internal/models"
	"github.com/Ijne/homepage_app/internal/tools"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func initDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error with .env in initDB: %s", err)
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_SSLMODE := os.Getenv("DB_SSLMODE")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Error parsing database config: %s", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Error creating database pool: %s", err)
	}

	dbPool = pool
	log.Println("Database connection pool initialized successfully")
}

func AddUser(ctx context.Context, username, email, password string) (int32, error) {
	once.Do(initDB)

	if err := dbPool.QueryRow(ctx, "SELECT 1 FROM users WHERE email = $1", email).Scan(new(int)); err == nil {
		return 0, fmt.Errorf("email %s already exists", email)
	}

	if !tools.ValidateEmail(email) {
		return 0, fmt.Errorf("invalid email format: %s", email)
	}

	passwordHash, err := tools.PasswordToHash(password)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %w", err)
	}

	var id int32
	if err := dbPool.QueryRow(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", username, email, passwordHash).Scan(&id); err != nil {
		return 0, fmt.Errorf("error adding user: %w", err)
	}

	return id, nil
}

func GetUser(ctx context.Context, id int32) (*models.User, error) {
	once.Do(initDB)

	var user models.User
	if err := dbPool.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &user, nil
}

func GetUserPassword(ctx context.Context, id int32) (string, error) {
	once.Do(initDB)

	var password string
	if err := dbPool.QueryRow(ctx, "SELECT password FROM users WHERE id = $1", id).Scan(&password); err != nil {
		return "", fmt.Errorf("error retrieving user password: %w", err)
	}

	return password, nil
}

func GetUserID(ctx context.Context, email string) (int32, error) {
	once.Do(initDB)

	var id int32
	if err := dbPool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", email).Scan(&id); err != nil {
		return 0, fmt.Errorf("error retrieving user ID: %w", err)
	}

	return id, nil
}

func GetUsername(ctx context.Context, id int32) (string, error) {
	once.Do(initDB)

	var name string
	if err := dbPool.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", id).Scan(&name); err != nil {
		return "", fmt.Errorf("error retrieving name: %w", err)
	}

	return name, nil
}
