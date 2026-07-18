package service

import (
	"context"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) List(ctx context.Context) ([]dto.AuditLogResponse, error) {
	entries, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.AuditLogResponse, 0, len(entries))
	for _, e := range entries {
		result = append(result, dto.AuditLogResponse{
			ID:        e.ID,
			UserID:    e.UserID,
			UserName:  e.UserName,
			Action:    e.Action,
			Entity:    e.Entity,
			EntityID:  e.EntityID,
			CreatedAt: e.CreatedAt,
		})
	}
	return result, nil
}
