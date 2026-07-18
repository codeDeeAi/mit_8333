package repository

import (
	"context"
	"database/sql"

	"UMSRMS/internal/models"

	squirrel "github.com/Masterminds/squirrel"
)

// UserWithRole is a user joined with its role name.
type UserWithRole struct {
	models.User
	RoleName string
}

// UserRepository handles users table operations.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := SQL.Insert("users").
		Columns("full_name", "email", "password_hash", "role_id", "phone").
		Values(user.FullName, user.Email, user.PasswordHash, user.RoleID, user.Phone).
		Suffix("RETURNING id, full_name, email, password_hash, role_id, phone, created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	created := &models.User{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&created.ID,
		&created.FullName,
		&created.Email,
		&created.PasswordHash,
		&created.RoleID,
		&created.Phone,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := SQL.Select(
		"id",
		"full_name",
		"email",
		"password_hash",
		"role_id",
		"phone",
		"created_at",
		"updated_at",
	).From("users").Where("email = ?", email).Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) List(ctx context.Context, roleName string) ([]UserWithRole, error) {
	query := SQL.Select(
		"u.id",
		"u.full_name",
		"u.email",
		"u.password_hash",
		"u.role_id",
		"u.phone",
		"u.created_at",
		"u.updated_at",
		"r.name AS role_name",
	).From("users u").Join("roles r ON r.id = u.role_id").OrderBy("u.created_at DESC")

	if roleName != "" {
		query = query.Where("r.name = ?", roleName)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]UserWithRole, 0)
	for rows.Next() {
		var user UserWithRole
		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.RoleID,
			&user.Phone,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.RoleName,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) UpdateRole(ctx context.Context, id, roleID int64) error {
	query := SQL.Update("users").
		Set("role_id", roleID).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where("id = ?", id)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := SQL.Delete("users").Where("id = ?", id)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := SQL.Select(
		"id",
		"full_name",
		"email",
		"password_hash",
		"role_id",
		"phone",
		"created_at",
		"updated_at",
	).From("users").Where("id = ?", id).Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
