package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
}
