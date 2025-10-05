package models

import "time"

type User struct {
	ID               int32  `json:"id"`
	Username         string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	SubscribersCount int    `json:"subscribers"`
}

type News struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Owner     int32     `json:"owner"`
	OwnerNAME string    `json:"owner_name"`
	CreatedAt time.Time `json:"created_at"`
	Image     string    `json:"img"`
}

type Subscription struct {
	UserID           int32  `json:"id"`
	Username         string `json:"username"`
	SubscribersCount int    `json:"subsribers"`
}
