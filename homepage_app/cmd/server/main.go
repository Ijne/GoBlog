package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"

	"github.com/Ijne/homepage_app/internal/handlers/homepage"
	"github.com/Ijne/homepage_app/internal/middlewares"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(middlewares.CheckAuth).Handle("/", http.HandlerFunc(homepage.HomepageHandler))
	r.Handle("/registration", http.HandlerFunc(homepage.RegisterHandler))
	r.Handle("/login", http.HandlerFunc(homepage.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(homepage.LogoutHandler))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
