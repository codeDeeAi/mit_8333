package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"UMSRMS/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header format")
	ErrMissingJWTSecret           = errors.New("jwt secret is required")
)

// UserClaims represents JWT claims used by the app.
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// JWTManager handles token generation and validation.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

// NewJWTManager creates a JWT manager using environment config.
func NewJWTManager(cfg *config.EnvConfig) (*JWTManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return nil, ErrMissingJWTSecret
	}
	if cfg.JWTExpireHours <= 0 {
		cfg.JWTExpireHours = 24
	}

	return &JWTManager{
		secret: []byte(cfg.JWTSecret),
		ttl:    time.Duration(cfg.JWTExpireHours) * time.Hour,
		issuer: cfg.AppName,
	}, nil
}

// GenerateToken creates a signed JWT for the provided user identity.
func (m *JWTManager) GenerateToken(userID, email, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)

	claims := UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// ParseToken validates a JWT and returns app claims.
func (m *JWTManager) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractBearerToken extracts JWT token from Authorization header.
func ExtractBearerToken(authHeader string) (string, error) {
	header := strings.TrimSpace(authHeader)
	if header == "" {
		return "", ErrMissingAuthorizationHeader
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrInvalidAuthorizationHeader
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrInvalidAuthorizationHeader
	}

	return token, nil
}
