package models

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
}
