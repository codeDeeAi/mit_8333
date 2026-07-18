package repository

import (
	"context"
	"database/sql"

	"UMSRMS/internal/models"
)

// AuditLogWithUser is an audit entry joined with the acting user's name.
type AuditLogWithUser struct {
	models.AuditLog
	UserName string
}

// AuditRepository handles audit_logs persistence.
type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Record(ctx context.Context, userID int64, action, entity, entityID string) error {
	var user any
	if userID > 0 {
		user = userID
	}

	query, args, err := SQL.Insert("audit_logs").
		Columns("user_id", "action", "entity", "entity_id").
		Values(user, action, entity, entityID).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *AuditRepository) List(ctx context.Context) ([]AuditLogWithUser, error) {
	query, args, err := SQL.Select(
		"a.id",
		"a.user_id",
		"a.action",
		"a.entity",
		"a.entity_id",
		"a.created_at",
		"COALESCE(u.full_name, 'System') AS user_name",
	).From("audit_logs a").
		LeftJoin("users u ON u.id = a.user_id").
		OrderBy("a.created_at DESC").
		Limit(200).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]AuditLogWithUser, 0)
	for rows.Next() {
		var entry AuditLogWithUser
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Action,
			&entry.Entity,
			&entry.EntityID,
			&entry.CreatedAt,
			&entry.UserName,
		); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}

	return logs, rows.Err()
}
