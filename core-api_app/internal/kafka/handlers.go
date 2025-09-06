package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Ijne/core-api_app/internal/models"
	"github.com/Ijne/core-api_app/internal/storage"
	"github.com/Ijne/core-api_app/internal/websockets"
)

func (c *Consumer) SetupHandlers() {
	c.RegisterHandler("notifications-out", func(msg *sarama.ConsumerMessage) error {
		var event struct {
			EventType string                 `json:"event_type"`
			Body      map[string]interface{} `json:"body"`
		}

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return err
		}

		fmt.Println(event, "From handlers")

		body := event.Body
		xuid := body["id_from"].(float64)
		id := int32(xuid)
		u, _ := storage.Get(context.Background(), id, "user")
		user := u.(models.User)
		switch event.EventType {
		case "subscribe":
			data := map[string]interface{}{
				"title":   body["title"].(string),
				"message": fmt.Sprintf(body["message"].(string), user.Username),
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Ошибка при создании JSON: %v", err)
				return err
			}

			xuid := body["id_to"].(float64)
			id := int32(xuid)
			websockets.WS_server.Broadcast(string(jsonData), id)
			fmt.Println("Sent from handlers type subscribe")
		case "new_post":
			data := map[string]interface{}{
				"title":   fmt.Sprintf(body["title"].(string), user.Username),
				"message": body["message"].(string),
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("Ошибка при создании JSON: %v", err)
				return err
			}

			xuid := body["id_from"].(float64)
			id := int32(xuid)
			for _, i := range storage.GetSubscribersID(context.Background(), id) {
				websockets.WS_server.Broadcast(string(jsonData), int32(i))
			}
		}

		return nil
	})
}
