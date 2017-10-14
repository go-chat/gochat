package model

import "time"

type GroupUser struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
