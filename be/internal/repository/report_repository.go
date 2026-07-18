package repository

import (
	"context"
	"database/sql"
)

// CategoryCount is a request count grouped by category name.
type CategoryCount struct {
	Name  string
	Count int64
}

// ReportRepository handles reporting aggregate queries.
type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Total(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM service_requests").Scan(&total)
	return total, err
}

func (r *ReportRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM service_requests GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func (r *ReportRepository) CountByCategory(ctx context.Context) ([]CategoryCount, error) {
	const query = `SELECT rc.name, COUNT(sr.id)
		FROM request_categories rc
		LEFT JOIN service_requests sr ON sr.category_id = rc.id
		GROUP BY rc.name
		ORDER BY rc.name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]CategoryCount, 0)
	for rows.Next() {
		var c CategoryCount
		if err := rows.Scan(&c.Name, &c.Count); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
