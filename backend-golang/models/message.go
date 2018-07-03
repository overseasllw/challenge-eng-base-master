package models

import (
	"time"
)

type (
	Message struct {
		MessageID      int64     `json:"message_id"`
		UUID           string    `json:"uuid"`
		MessageType    string    `json:"message_type"`
		CreatedAt      time.Time `json:"created_at,omitempty"`
		MessageContent *string   `json:"message"`
		Room           *string   `json:"room"`
		RoomId         int64     `room_id`
		User
	}
)
