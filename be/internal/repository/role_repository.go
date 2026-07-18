package repository

import (
	"context"
	"database/sql"

	"UMSRMS/internal/models"
)

// RoleRepository handles roles table operations and seeding.
type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := SQL.Select("id", "name", "description", "created_at", "updated_at").
		From("roles").
		Where("name = ?", name).
		Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	role := &models.Role{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	query := SQL.Select("id", "name", "description", "created_at", "updated_at").
		From("roles").
		Where("id = ?", id).
		Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	role := &models.Role{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]models.Role, error) {
	query, args, err := SQL.Select("id", "name", "description", "created_at", "updated_at").
		From("roles").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *RoleRepository) SeedDefaultRoles(ctx context.Context) (int64, error) {
	query, args, err := SQL.Insert("roles").
		Columns("name", "description").
		Values(models.RoleStudentStaff, "Student or staff user").
		Values(models.RoleMaintenanceOfficer, "Maintenance officer").
		Values(models.RoleAdmin, "System administrator").
		Suffix("ON CONFLICT (name) DO NOTHING").
		ToSql()
	if err != nil {
		return 0, err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
