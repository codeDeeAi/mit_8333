package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/models"
	"UMSRMS/internal/repository"
	"UMSRMS/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotFound       = errors.New("default role not found")
)

// AuthService handles registration and login flows.
type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	jwt      *utils.JWTManager
	banList  *utils.TokenBanList
}

func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	jwtManager *utils.JWTManager,
	banList *utils.TokenBanList,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		jwt:      jwtManager,
		banList:  banList,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, err := s.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	role, err := s.resolveRegistrationRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:     req.FullName,
		Email:        strings.ToLower(req.Email),
		PasswordHash: string(hash),
		RoleID:       role.ID,
		Phone:        req.Phone,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, expiresAt, err := s.jwt.GenerateToken(strconv.FormatInt(createdUser.ID, 10), createdUser.Email, role.Name)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        createdUser.ID,
			FullName:  createdUser.FullName,
			Email:     createdUser.Email,
			RoleID:    createdUser.RoleID,
			Role:      role.Name,
			Phone:     createdUser.Phone,
			CreatedAt: createdUser.CreatedAt,
			UpdatedAt: createdUser.UpdatedAt,
		},
	}, nil
}

// resolveRegistrationRole returns the chosen role, or the default student/staff
// role when none is supplied.
func (s *AuthService) resolveRegistrationRole(ctx context.Context, roleID int64) (*models.Role, error) {
	if roleID > 0 {
		role, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrRoleNotFound
			}
			return nil, err
		}
		return role, nil
	}

	role, err := s.roleRepo.GetByName(ctx, models.RoleStudentStaff)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	roleName := ""
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err == nil {
		roleName = role.Name
	}

	token, expiresAt, err := s.jwt.GenerateToken(strconv.FormatInt(user.ID, 10), user.Email, roleName)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			RoleID:    user.RoleID,
			Role:      roleName,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	roleName := ""
	if role, roleErr := s.roleRepo.GetByID(ctx, user.RoleID); roleErr == nil {
		roleName = role.Name
	}

	return &dto.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		RoleID:    user.RoleID,
		Role:      roleName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// RegistrationData returns reference data the client needs for the sign-up
// screen, currently the list of available roles.
func (s *AuthService) RegistrationData(ctx context.Context) (*dto.RegistrationDataResponse, error) {
	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := &dto.RegistrationDataResponse{Roles: make([]dto.RoleResponse, 0, len(roles))}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return resp, nil
}

// RevokeToken bans a token until it expires, logging the user out. The caller
// (via RequireAuth middleware) has already validated the token.
func (s *AuthService) RevokeToken(token string, expiresAt time.Time) {
	if s.banList == nil || strings.TrimSpace(token) == "" {
		return
	}

	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().Add(24 * time.Hour)
	}

	s.banList.Ban(token, expiresAt)
}
