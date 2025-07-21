package middlewares

import (
	"net/http"

	"github.com/Ijne/core-api_app/internal/models"
	"github.com/Ijne/core-api_app/internal/tools"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := tools.GetUserClaimsFromCookie(r)
		if err != nil {
			action := r.URL.Query().Get("action")
			var data = struct {
				User models.User
				News []models.News
			}{
				User: models.User{ID: 0},
				News: []models.News{},
			}
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

		next.ServeHTTP(w, r)
	})
}
