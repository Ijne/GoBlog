package models

import "time"

type User struct {
	ID          int32  `json:"id"`
	Username    string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Subscribers int    `json:"subscribers"`
}

type News struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Owner     int32     `json:"owner"`
	OwnerNAME string    `json:"owner_name"`
	CreatedAt time.Time `json:"created_at"`
}
