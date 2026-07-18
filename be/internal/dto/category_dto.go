package dto

// CategoryResponse is a request category exposed to clients.
type CategoryResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}
