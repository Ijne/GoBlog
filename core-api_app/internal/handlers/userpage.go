package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ijne/core-api_app/internal/models"
	"github.com/Ijne/core-api_app/internal/storage"
	"github.com/Ijne/core-api_app/internal/tools"
)

func UserPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idString := r.URL.Query().Get("id")
		id, _ := strconv.ParseInt(idString, 10, 32)
		user, err := storage.Get(context.Background(), int32(id), "user")
		if err != nil {
			log.Fatal("Can't read user id")
		}

		var news interface{}
		news, err = storage.Get(context.Background(), int32(id), "news")
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

		var guest models.User
		guest, err = tools.GetUserClaimsFromCookie(r)
		if err != nil {
			log.Println(err)
			return
		}
		user_c := user.(models.User)
		var data = struct {
			User         models.User
			News         []models.News
			IsSubscribed bool
		}{
			User:         user_c,
			News:         newsSlice,
			IsSubscribed: storage.GetSubscription(context.Background(), guest.ID, user_c.ID),
		}
		log.Println(storage.GetSubscription(context.Background(), guest.ID, user_c.ID))

		tools.RenderTemplate(w, "userpage.html", data)

	default:
		log.Println("Method not allowed(UserPageHandler)", r.Method)
		return
	}
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		user, err := tools.GetUserClaimsFromCookie(r)
		if err != nil {
			log.Println(err)
			return
		}

		var data struct {
			TargetID int32  `json:"target_id"`
			Action   string `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Fatal("Can't decode response FROM subscribe")
		}
		switch data.Action {
		case "subscribe":
			if err := storage.AddSubscription(context.Background(), user.ID, data.TargetID); err != nil {
				log.Fatal(err)
			}
			var data = struct{}{}
			json.NewEncoder(w).Encode(data)
			log.Println("Sucssesfully subscribed!")
		case "unsubscribe":
			if err := storage.DelSubscription(context.Background(), user.ID, data.TargetID); err != nil {
				log.Fatal(err)
			}
			var data = struct{}{}
			json.NewEncoder(w).Encode(data)
			log.Println("Sucssesfully unsubscribed!")
		default:
			log.Println("Unrecognized action(SubscribeHandler)")
		}
	default:
		log.Println("Method not allowed(SubscribeHandler)")
		return
	}
}
