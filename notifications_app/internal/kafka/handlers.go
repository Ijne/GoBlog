package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

func (c *Consumer) SetupHandlers() {
	c.RegisterHandler("notifications", func(msg *sarama.ConsumerMessage) error {
		var event struct {
			EventType string                 `json:"event_type"`
			Body      map[string]interface{} `json:"body"`
		}

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return err
		}

		var title, message string
		switch event.EventType {
		case "subscribe":
			title = "У Вас новый подписчик!"
			message = "На Вас подписался %s"
			SendMessage(map[string]interface{}{
				"event_type": "subscribe",
				"body": map[string]interface{}{
					"title":   title,
					"message": message,
					"id_from": event.Body["from"],
					"id_to":   event.Body["to"],
				},
			})
		case "new_post":
			title = "У %s новый пост"
			message = "Скорее смотрите!"
			SendMessage(map[string]interface{}{
				"event_type": "new_post",
				"body": map[string]interface{}{
					"title":   title,
					"message": message,
					"id_from": event.Body["from"],
				},
			})
		}

		fmt.Println(event)

		return nil
	})
}
