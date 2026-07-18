package middleware

import (
	"strconv"
	"time"

	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID         = "auth_user_id"
	ContextUserEmail      = "auth_user_email"
	ContextUserRole       = "auth_user_role"
	ContextToken          = "auth_token"
	ContextTokenExpiresAt = "auth_token_expires_at"
)

// RequireAuth validates the bearer token, rejects revoked tokens, and stores the
// authenticated user's identity and token in the request context.
func RequireAuth(jwtManager *utils.JWTManager, banList *utils.TokenBanList) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			utils.Unauthenticated(c, "Authorization header is required")
			c.Abort()
			return
		}

		if banList != nil && banList.IsBanned(token) {
			utils.Unauthenticated(c, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			utils.Unauthenticated(c, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(claims.UserID, 10, 64)
		if err != nil {
			utils.Unauthenticated(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, userID)
		c.Set(ContextUserEmail, claims.Email)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextToken, token)
		if claims.ExpiresAt != nil {
			c.Set(ContextTokenExpiresAt, claims.ExpiresAt.Time.UTC())
		}

		c.Next()
	}
}

// UserID returns the authenticated user's ID from the context.
func UserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ContextUserID)
	id, valid := v.(int64)
	return id, ok && valid
}

// UserRole returns the authenticated user's role from the context.
func UserRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(ContextUserRole)
	role, valid := v.(string)
	return role, ok && valid && role != ""
}

// RequireRoles allows only the provided roles to access a route.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, ok := UserRole(c)
		if !ok {
			utils.Unauthorized(c, "Authenticated user role is missing")
			c.Abort()
			return
		}

		if _, exists := allowed[role]; !exists {
			utils.Unauthorized(c, "You are not allowed to perform this action")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Token returns the raw bearer token from the context.
func Token(c *gin.Context) string {
	v, _ := c.Get(ContextToken)
	token, _ := v.(string)
	return token
}

// TokenExpiresAt returns the token's expiry from the context.
func TokenExpiresAt(c *gin.Context) (time.Time, bool) {
	v, ok := c.Get(ContextTokenExpiresAt)
	t, valid := v.(time.Time)
	return t, ok && valid
}
