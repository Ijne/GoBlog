package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Ijne/core-api_app/internal/models"
	"github.com/Ijne/core-api_app/internal/storage"
	"github.com/Ijne/core-api_app/internal/tools"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user, err := tools.GetUserClaimsFromCookie(r)
		if err != nil {
			log.Println("Error getting user claims:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var news interface{}
		news, err = storage.Get(context.Background(), 0, "news")
		if err != nil {
			log.Println(err)
			return
		}
		newsSlice, ok := news.([]models.News)
		if !ok {
			log.Println(err)
			return
		}
		if len(newsSlice) == 1 && newsSlice[0].ID == 0 {
			news = []models.News{}
		}

		var data = struct {
			User models.User
			News []models.News
		}{
			User: user,
			News: newsSlice,
		}

		tools.RenderTemplate(w, "homepage.html", data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
