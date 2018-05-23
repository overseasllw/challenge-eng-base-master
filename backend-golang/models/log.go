package models

import "time"

type UserActivityHistory struct {
	ActivityID int64     `json:"activity_id"`
	UserID     int64     `json:"user_id"`
	ActiveAt   time.Time `json:"active_at"`
}
