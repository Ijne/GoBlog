package tools

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

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
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, c := range email {
		if c == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if c == '.' && (i == 0 || i == len(email)-1 || email[i-1] == '@') {
			return false
		}
	}
	return at == 1
}

func ExtractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
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
