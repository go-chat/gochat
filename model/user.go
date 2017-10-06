package model

import "time"

// User ..
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Avatar    string    `json:"avatar"`
	Meta      string    `json:"meta"`
	Salt      string    `json:"salt"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
