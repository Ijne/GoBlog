package middlewares

import (
	"log"
	"net/http"

	"github.com/Ijne/homepage_app/internal/models"
	"github.com/Ijne/homepage_app/internal/tools"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tools.ExtractTokenFromCookie(r)

		if token == "" {
			action := r.URL.Query().Get("action")
			data := models.User{ID: 0, Username: "", Email: ""}
			switch action {
			case "login":
				tools.RenderTemplate(w, "login.html", data)
			case "register":
				tools.RenderTemplate(w, "register.html", data)
			default:
				tools.RenderTemplate(w, "register.html", data)
			}
			return
		}

		claims, err := tools.ValidateToken(token)
		if err != nil || claims == nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
