package dto

import "time"

// NotificationResponse is an in-app notification returned to the owning user.
type NotificationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
