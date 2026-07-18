package service

import (
	"context"
	"database/sql"
	"errors"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/repository"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) List(ctx context.Context, userID int64) ([]dto.NotificationResponse, error) {
	notifications, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, dto.NotificationResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			Message:   n.Message,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}
	return result, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, id, userID int64) error {
	if err := s.repo.MarkRead(ctx, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotificationNotFound
		}
		return err
	}
	return nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllRead(ctx, userID)
}
