package models

type User struct {
	ID       int32  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"  `
	Password string `json:"password" db:"password"`
}
