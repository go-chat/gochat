package model

import "time"

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
