package dto

// CategoryCountResponse is a request count for one category.
type CategoryCountResponse struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// ReportSummaryResponse is the admin dashboard summary.
type ReportSummaryResponse struct {
	Total      int64                   `json:"total"`
	ByStatus   map[string]int64        `json:"by_status"`
	ByCategory []CategoryCountResponse `json:"by_category"`
}
