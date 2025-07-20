package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ijne/homepage_app/internal/models"
	"github.com/Ijne/homepage_app/internal/storage"
	"github.com/Ijne/homepage_app/internal/tools"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data struct {
			Username        string `json:"username"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirm_password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if data.Password != data.ConfirmPassword {
			log.Println("Passwords do not match")
			http.Error(w, "Passwords do not match", http.StatusForbidden)
			return
		}

		password, err := tools.PasswordToHash(data.Password)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}
		id, err := storage.Add(context.Background(), models.User{
			Username: data.Username,
			Email:    data.Email,
			Password: password,
		})
		if err != nil {
			log.Println("Error adding user:", err)
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		tools.MakeCookieAfterLogin(w, id, data.Username, data.Email)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := storage.GetUserByEmail(context.Background(), data.Email)
		if err != nil {
			log.Println("Error retrieving user:", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		tools.MakeCookieAfterLogin(w, user.ID, user.Username, data.Email)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
