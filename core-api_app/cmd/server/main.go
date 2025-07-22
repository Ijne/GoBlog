package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"

	"github.com/Ijne/core-api_app/internal/handlers/auth"
	"github.com/Ijne/core-api_app/internal/handlers/homepage"
	"github.com/Ijne/core-api_app/internal/handlers/profile"
	"github.com/Ijne/core-api_app/internal/middlewares"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(middlewares.CheckAuth).Handle("/", http.HandlerFunc(homepage.HomepageHandler))
	r.Handle("/registration", http.HandlerFunc(auth.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(auth.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(auth.LogoutHandler))
	r.With(middlewares.CheckAuth).Handle("/profile", http.HandlerFunc(profile.ProfileHandler))
	r.With(middlewares.CheckAuth).Handle("/news", http.HandlerFunc(profile.NewsHandler))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
