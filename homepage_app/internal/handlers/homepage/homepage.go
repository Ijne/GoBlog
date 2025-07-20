package homepage

import (
	"log"
	"net/http"

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
