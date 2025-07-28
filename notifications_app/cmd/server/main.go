package main

import (
	"fmt"
	"net/http"

	"github.com/Ijne/notifications_app/internal/kafka"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	fmt.Println("1")

	consumer, err := kafka.NewConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.SetupHandlers()
	go consumer.Start()

	fmt.Println("2")

	if err := http.ListenAndServe("0.0.0.0:8070", r); err != nil {
		panic(err)
	}
}
