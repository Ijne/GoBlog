package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/websocket"

	"github.com/Ijne/core-api_app/internal/handlers"
	"github.com/Ijne/core-api_app/internal/kafka"
	"github.com/Ijne/core-api_app/internal/middlewares"
	"github.com/Ijne/core-api_app/internal/websockets"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")
	KAFKA_HOST := os.Getenv("KAFKA_HOST")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(middlewares.CheckAuth).Handle("/ws", websocket.Handler(websockets.WS_server.HandleWS))

	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "internal/static")
	fmt.Println(staticDir)
	fs := http.FileServer(http.Dir(staticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.With(middlewares.CheckAuth).Handle("/", http.HandlerFunc(handlers.HomepageHandler))
	r.Handle("/registration", http.HandlerFunc(handlers.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(handlers.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))
	r.With(middlewares.CheckAuth).Handle("/profile", http.HandlerFunc(handlers.ProfileHandler))
	r.With(middlewares.CheckAuth).Handle("/news", http.HandlerFunc(handlers.NewsHandler))
	r.With(middlewares.CheckAuth).Handle("/user", http.HandlerFunc(handlers.UserPageHandler))
	r.With(middlewares.CheckAuth).Handle("/subscribe", http.HandlerFunc(handlers.SubscribeHandler))

	consumer, err := kafka.NewConsumer([]string{fmt.Sprintf("%s:9092", KAFKA_HOST)})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.SetupHandlers()
	go consumer.Start()

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", PORT), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
