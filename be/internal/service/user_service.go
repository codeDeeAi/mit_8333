package service

import (
	"context"
	"database/sql"
	"errors"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/repository"
)

var (
	ErrCannotModifySelf = errors.New("you cannot modify your own account")
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserService) List(ctx context.Context, roleName string) ([]dto.UserResponse, error) {
	users, err := s.userRepo.List(ctx, roleName)
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, toUserResponse(user))
	}
	return result, nil
}

func (s *UserService) UpdateRole(ctx context.Context, actorID, userID int64, roleName string) (*dto.UserResponse, error) {
	if actorID == userID {
		return nil, ErrCannotModifySelf
	}

	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	if err := s.userRepo.UpdateRole(ctx, userID, role.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	updated, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := dto.UserResponse{
		ID:        updated.ID,
		FullName:  updated.FullName,
		Email:     updated.Email,
		RoleID:    updated.RoleID,
		Role:      role.Name,
		Phone:     updated.Phone,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}
	return &response, nil
}

func (s *UserService) Delete(ctx context.Context, actorID, userID int64) error {
	if actorID == userID {
		return ErrCannotModifySelf
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func toUserResponse(user repository.UserWithRole) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		RoleID:    user.RoleID,
		Role:      user.RoleName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
