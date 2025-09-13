package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ijne/notifications_app/internal/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	time.Sleep(30 * time.Second)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := os.Getenv("PORT")
	KAFKA_HOST := os.Getenv("KAFKA_HOST")

	r := chi.NewRouter()

	fmt.Println("1")

	consumer, err := kafka.NewConsumer([]string{fmt.Sprintf("%s:9092", KAFKA_HOST)})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	consumer.SetupHandlers()
	go consumer.Start()

	fmt.Println("2")

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", PORT), r); err != nil {
		panic(err)
	}
}
