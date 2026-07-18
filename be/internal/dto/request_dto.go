package dto

import "time"

// CreateServiceRequestRequest is the multipart/form-data payload for creating a request.
type CreateServiceRequestRequest struct {
	Title       string `form:"title" binding:"required,min=3,max=200"`
	Description string `form:"description" binding:"required,min=5"`
	CategoryID  int64  `form:"category_id" binding:"required,gt=0"`
	Location    string `form:"location" binding:"required,min=2,max=255"`
	Priority    string `form:"priority" binding:"omitempty,oneof=low medium high"`
}

// UpdateServiceRequestStatusRequest is the payload for request status updates.
type UpdateServiceRequestStatusRequest struct {
	Status string  `json:"status" binding:"required,oneof=pending assigned in_progress completed rejected"`
	Note   *string `json:"note,omitempty" binding:"omitempty,max=1000"`
}

// ListServiceRequestsQuery captures supported list filters.
type ListServiceRequestsQuery struct {
	Q          string `form:"q"`
	Status     string `form:"status" binding:"omitempty,oneof=pending assigned in_progress completed rejected"`
	CategoryID int64  `form:"category_id" binding:"omitempty,gt=0"`
	Priority   string `form:"priority" binding:"omitempty,oneof=low medium high"`
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Sort       string `form:"sort"`
}

// ServiceRequestResponse is the public request payload returned by the API.
type ServiceRequestResponse struct {
	ID                  int64     `json:"id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	CategoryID          int64     `json:"category_id"`
	CategoryName        string    `json:"category_name,omitempty"`
	CreatedBy           int64     `json:"created_by"`
	CreatedByName       string    `json:"created_by_name,omitempty"`
	Location            string    `json:"location"`
	Priority            string    `json:"priority"`
	Status              string    `json:"status"`
	EvidencePath        *string   `json:"evidence_path,omitempty"`
	AssignedOfficerID   *int64    `json:"assigned_officer_id,omitempty"`
	AssignedOfficerName *string   `json:"assigned_officer_name,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// AssignServiceRequestRequest is the payload for assigning a request to an officer.
type AssignServiceRequestRequest struct {
	OfficerID int64 `json:"officer_id" binding:"required,gt=0"`
}

// StatusLogResponse represents a request status change log entry.
type StatusLogResponse struct {
	ID        int64     `json:"id"`
	RequestID int64     `json:"request_id"`
	ChangedBy int64     `json:"changed_by"`
	OldStatus *string   `json:"old_status,omitempty"`
	NewStatus string    `json:"new_status"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ServiceRequestDetailResponse returns a request with its status history.
type ServiceRequestDetailResponse struct {
	Request    ServiceRequestResponse `json:"request"`
	StatusLogs []StatusLogResponse    `json:"status_logs"`
}

// PaginationMeta returns list paging metadata.
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ServiceRequestListResponse wraps the paginated request list.
type ServiceRequestListResponse struct {
	Items      []ServiceRequestResponse `json:"items"`
	Pagination PaginationMeta           `json:"pagination"`
}
