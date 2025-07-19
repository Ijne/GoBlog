package tools

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ijne/homepage_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func PasswordToHash(password string) (string, error) {
	bcryptCost := 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(password, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	// Проверка длины email
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Проверка наличия ровно одного символа @
	atCount := 0
	atPos := -1
	for i, c := range email {
		if c == '@' {
			atCount++
			atPos = i
			if atCount > 1 {
				return false
			}
		}
	}
	if atCount != 1 {
		return false
	}

	// Разделение на локальную часть и домен
	localPart := email[:atPos]
	domain := email[atPos+1:]

	// Проверка локальной части
	if len(localPart) < 1 || len(localPart) > 64 {
		return false
	}
	if localPart[0] == '.' || localPart[len(localPart)-1] == '.' {
		return false
	}
	for i := 0; i < len(localPart)-1; i++ {
		if localPart[i] == '.' && localPart[i+1] == '.' {
			return false
		}
	}

	// Проверка домена
	if len(domain) < 1 || len(domain) > 253 {
		return false
	}
	if domain[0] == '.' || domain[len(domain)-1] == '.' {
		return false
	}
	for i := 0; i < len(domain)-1; i++ {
		if domain[i] == '.' && domain[i+1] == '.' {
			return false
		}
	}

	// Проверка допустимых символов (упрощенная)
	for _, c := range localPart {
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') &&
			!(c >= '0' && c <= '9') && c != '.' && c != '_' && c != '-' && c != '+' {
			return false
		}
	}

	for _, c := range domain {
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') &&
			!(c >= '0' && c <= '9') && c != '.' && c != '-' {
			return false
		}
	}

	return true
}

func ExtractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func MakeCookieAfterLogin(w http.ResponseWriter, id int32, username, email string) {
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "Error loading environment variables", http.StatusInternalServerError)
		return
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"redirect": "/",
	})

}

func GetCookieClaims(r *http.Request) (models.User, error) {
	token := ExtractTokenFromCookie(r)
	claims, err := ValidateToken(token)
	if err != nil || claims == nil {
		log.Println(err)
		return models.User{}, fmt.Errorf("unauthorized")
	}

	id, i_ok := claims.(jwt.MapClaims)["id"].(float64)
	username, u_ok := claims.(jwt.MapClaims)["username"].(string)
	email, e_ok := claims.(jwt.MapClaims)["email"].(string)

	if !u_ok || !e_ok || !i_ok {
		log.Println("Invalid token claims", u_ok, e_ok, i_ok)
		return models.User{}, fmt.Errorf("invalid token claims")
	}

	user := models.User{
		ID:       int32(id),
		Username: username,
		Email:    email,
	}
	return user, nil
}

func ValidateToken(inputToken string) (jwt.Claims, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected algorithm: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	templates := template.Must(template.ParseFiles(
		"internal/templates/base.html",
		"internal/templates/"+tmpl,
	))

	err := templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error with execute template")
	}
	return nil
}
