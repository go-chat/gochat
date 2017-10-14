package model

type Message struct {
	ID       int64  `json:"id"`
	GroupID  int64  `json:"group_id"`
	SenderID int64  `json:"sender_id"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
}

const (
	MessageTypeText = iota
	MessageTypeImage
)
