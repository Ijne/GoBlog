package homepage

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ijne/homepage_app/internal/storage"
	"github.com/Ijne/homepage_app/internal/tools"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user, err := tools.GetCookieClaims(r)
		if err != nil {
			log.Println("Error getting user claims:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tools.RenderTemplate(w, "homepage.html", user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

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

		id, err := storage.AddUser(context.Background(), data.Username, data.Email, data.Password)
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

		userID, err := storage.GetUserID(context.Background(), data.Email)
		if err != nil {
			log.Println("Error retrieving user ID:", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		username, err := storage.GetUsername(context.Background(), userID)
		if err != nil {
			log.Println("Error retrieving user ID:", err)
			http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
			return
		}
		hashedPassword, err := storage.GetUserPassword(context.Background(), userID)
		if err != nil {
			log.Println("Error retrieving user password:", err)
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}
		if !tools.ValidatePassword(data.Password, hashedPassword) {
			log.Println("Invalid email or password")
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		tools.MakeCookieAfterLogin(w, userID, username, data.Email)
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
