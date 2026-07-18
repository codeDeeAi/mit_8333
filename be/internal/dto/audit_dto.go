package dto

import "time"

// AuditLogResponse is an audit trail entry returned to admins.
type AuditLogResponse struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id,omitempty"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	CreatedAt time.Time `json:"created_at"`
}
