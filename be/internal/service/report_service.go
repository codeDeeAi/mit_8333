package service

import (
	"context"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/models"
	"UMSRMS/internal/repository"
)

type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) Summary(ctx context.Context) (*dto.ReportSummaryResponse, error) {
	total, err := s.repo.Total(ctx)
	if err != nil {
		return nil, err
	}

	rawStatus, err := s.repo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure every known status is present (defaulting to zero).
	byStatus := map[string]int64{
		models.StatusPending:    0,
		models.StatusAssigned:   0,
		models.StatusInProgress: 0,
		models.StatusCompleted:  0,
		models.StatusRejected:   0,
	}
	for status, count := range rawStatus {
		byStatus[status] = count
	}

	categories, err := s.repo.CountByCategory(ctx)
	if err != nil {
		return nil, err
	}

	byCategory := make([]dto.CategoryCountResponse, 0, len(categories))
	for _, c := range categories {
		byCategory = append(byCategory, dto.CategoryCountResponse{Name: c.Name, Count: c.Count})
	}

	return &dto.ReportSummaryResponse{
		Total:      total,
		ByStatus:   byStatus,
		ByCategory: byCategory,
	}, nil
}
