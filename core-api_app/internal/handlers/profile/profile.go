package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ijne/core-api_app/internal/models"
	"github.com/Ijne/core-api_app/internal/storage"
	"github.com/Ijne/core-api_app/internal/tools"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user, err := tools.GetUserClaimsFromCookie(r)
		if err != nil {
			log.Println("Error getting user claims:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var news interface{}
		news, err = storage.Get(context.Background(), user.ID, "news")
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

		tools.RenderTemplate(w, "profile.html", data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	switch action {
	case "add":
		var news models.News
		if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
			log.Println(err)
			return
		}
		user, err := tools.GetUserClaimsFromCookie(r)
		if err != nil {
			log.Println(err)
			return
		}
		news.Owner = user.ID
		news.CreatedAt = time.Now()

		id, err := storage.Add(context.Background(), news)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		var data = struct {
			ok bool
		}{
			ok: true,
		}
		json.NewEncoder(w).Encode(data)
		log.Println("Sucssesfully created news with id:", id)
	default:
		log.Println("Not allowed method")
		return
	}
}
