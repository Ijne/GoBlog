package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Ijne/core-api_app/internal/models"
	"github.com/jackc/pgx/v5"
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
	log.Println("Database connection established")
}

// USER TABLE METHODS
type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{dbPool: pool}
}

func (ur *UserRepository) Add(ctx context.Context, user models.User) (int32, error) {
	var id int32
	err := ur.dbPool.QueryRow(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}
	return id, nil
}

func (ur *UserRepository) Get(ctx context.Context, id int32) (*models.User, error) {
	var user models.User
	err := ur.dbPool.QueryRow(ctx, "SELECT id, username, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by id: %w", err)
	}
	return &user, nil
}

//

// NEWS TABLE METHODS
type NewsRepository struct {
	dbPool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{dbPool: pool}
}

func (nr *NewsRepository) Add(ctx context.Context, news models.News) (int32, error) {
	var id int32
	if err := dbPool.QueryRow(context.Background(), "INSERT INTO news (title, body, owner, created_at) VALUES ($1, $2, $3, $4) RETURNING id", news.Title, news.Body, news.Owner, news.CreatedAt).Scan(&id); err != nil {
		log.Println(err)
		return 0, fmt.Errorf("error inserting news: %w", err)
	}
	return id, nil
}

func (nr *NewsRepository) Get(ctx context.Context, id int32) ([]models.News, error) {
	var news []models.News
	var rows pgx.Rows
	var err error
	if id == 0 {
		rows, err = nr.dbPool.Query(ctx, "SELECT id, title, body, created_at FROM news ORDER BY created_at DESC")
	} else {
		rows, err = nr.dbPool.Query(ctx, "SELECT id, title, body, created_at FROM news WHERE owner = $1 ORDER BY created_at DESC", id)
	}
	if err != nil {
		return []models.News{}, fmt.Errorf("error fetching user by id: %w", err)
	}
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			log.Println("Error with scanning news:", err)
		}
		news = append(news, n)
	}
	fmt.Println("news:", news)
	return news, nil
}

func (nr *NewsRepository) Del(ctx context.Context, id int32) error {
	_, err := dbPool.Exec(ctx, "DELETE FROM news WHERE id=$1", id)
	return err
}

//

// MAIN ADD AND GET FUNCTIONS AND SPECIALS
func Add(ctx context.Context, data interface{}) (int32, error) {
	once.Do(initDB)
	user, ok := data.(models.User)
	if ok {
		return NewUserRepository(dbPool).Add(ctx, user)
	}
	news, ok := data.(models.News)
	if ok {
		return NewNewsRepository(dbPool).Add(ctx, news)
	}
	return 0, fmt.Errorf("invalid data type")
}

func Get(ctx context.Context, id int32, t string) (interface{}, error) {
	once.Do(initDB)
	switch t {
	case "user":
		return NewUserRepository(dbPool).Get(ctx, id)
	case "news":
		return NewNewsRepository(dbPool).Get(ctx, id)
	default:
		log.Println("Invalid type for Get method:", t)
		return nil, fmt.Errorf("invalid type, expected 'user'")
	}
}

func Del(ctx context.Context, id int32, t string) error {
	once.Do(initDB)
	switch t {
	case "news":
		return NewNewsRepository(dbPool).Del(ctx, id)
	default:
		return fmt.Errorf("method not allowed")
	}
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	once.Do(initDB)
	var user models.User
	err := dbPool.QueryRow(ctx, "SELECT id, username, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}
	return &user, nil
}

//
