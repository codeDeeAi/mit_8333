package service

import (
	"context"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return result, nil
}
