package repository

import (
	"context"
	"database/sql"

	"UMSRMS/internal/models"
)

// CategoryRepository handles request_categories table reads.
type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.RequestCategory, error) {
	query, args, err := SQL.Select("id", "name", "description", "created_at", "updated_at").
		From("request_categories").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.RequestCategory, 0)
	for rows.Next() {
		var category models.RequestCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}
