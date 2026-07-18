package repository

import (
	"context"
	"database/sql"

	"UMSRMS/internal/models"
)

// NotificationRepository handles notifications persistence.
type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, userID int64, message string) error {
	query, args, err := SQL.Insert("notifications").
		Columns("user_id", "message").
		Values(userID, message).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// CreateForRole inserts a notification for every user holding the given role.
func (r *NotificationRepository) CreateForRole(ctx context.Context, roleName, message string) error {
	const query = `INSERT INTO notifications (user_id, message)
		SELECT u.id, $1 FROM users u JOIN roles r ON r.id = u.role_id WHERE r.name = $2`
	_, err := r.db.ExecContext(ctx, query, message, roleName)
	return err
}

func (r *NotificationRepository) List(ctx context.Context, userID int64) ([]models.Notification, error) {
	query, args, err := SQL.Select("id", "user_id", "message", "is_read", "created_at").
		From("notifications").
		Where("user_id = ?", userID).
		OrderBy("created_at DESC").
		Limit(100).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]models.Notification, 0)
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID int64) error {
	query, args, err := SQL.Update("notifications").
		Set("is_read", true).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
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

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID int64) error {
	query, args, err := SQL.Update("notifications").
		Set("is_read", true).
		Where("user_id = ?", userID).
		Where("is_read = ?", false).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
