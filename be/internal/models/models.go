package models

import (
	"encoding/json"
	"time"
)

const (
	RoleStudentStaff       = "student_staff"
	RoleMaintenanceOfficer = "maintenance_officer"
	RoleAdmin              = "admin"
)

const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

const (
	StatusPending    = "pending"
	StatusAssigned   = "assigned"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusRejected   = "rejected"
)

// Role maps to the roles table.
type Role struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// User maps to the users table.
type User struct {
	ID           int64     `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	RoleID       int64     `json:"role_id" db:"role_id"`
	Phone        *string   `json:"phone,omitempty" db:"phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RequestCategory maps to the request_categories table.
type RequestCategory struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ServiceRequest maps to the service_requests table.
type ServiceRequest struct {
	ID           int64     `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	CategoryID   int64     `json:"category_id" db:"category_id"`
	CreatedBy    int64     `json:"created_by" db:"created_by"`
	Location     string    `json:"location" db:"location"`
	Priority     string    `json:"priority" db:"priority"`
	Status       string    `json:"status" db:"status"`
	EvidencePath *string   `json:"evidence_path,omitempty" db:"evidence_path"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Assignment maps to the assignments table.
type Assignment struct {
	ID         int64     `json:"id" db:"id"`
	RequestID  int64     `json:"request_id" db:"request_id"`
	OfficerID  int64     `json:"officer_id" db:"officer_id"`
	AssignedBy int64     `json:"assigned_by" db:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

// StatusLog maps to the status_logs table.
type StatusLog struct {
	ID        int64     `json:"id" db:"id"`
	RequestID int64     `json:"request_id" db:"request_id"`
	ChangedBy int64     `json:"changed_by" db:"changed_by"`
	OldStatus *string   `json:"old_status,omitempty" db:"old_status"`
	NewStatus string    `json:"new_status" db:"new_status"`
	Note      *string   `json:"note,omitempty" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Notification maps to the notifications table.
type Notification struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Message   string    `json:"message" db:"message"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuditLog maps to the audit_logs table.
type AuditLog struct {
	ID        int64           `json:"id" db:"id"`
	UserID    *int64          `json:"user_id,omitempty" db:"user_id"`
	Action    string          `json:"action" db:"action"`
	Entity    string          `json:"entity" db:"entity"`
	EntityID  string          `json:"entity_id" db:"entity_id"`
	Metadata  json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
