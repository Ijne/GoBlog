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

	"github.com/Ijne/core-api_app/internal/handlers"
	"github.com/Ijne/core-api_app/internal/middlewares"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
