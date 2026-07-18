package handler

import (
	"errors"
	"net/http"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/middleware"
	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler exposes auth endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// bindJSON binds and validates a request body, returning false (and writing an
// error response) when the payload is malformed or fails validation.
func bindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		if details, ok := utils.FormatValidationErrors(err); ok {
			utils.ValidationError(c, "Validation failed", details)
			return false
		}
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return false
	}
	return true
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a student/staff account and returns an auth token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterRequest  true  "Registration payload"
// @Success      201      {object}  utils.APIResponse{data=dto.AuthResponse}
// @Failure      409      {object}  utils.APIResponse
// @Failure      422      {object}  utils.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			utils.Conflict(c, "Email already exists", gin.H{"email": req.Email})
		case errors.Is(err, service.ErrRoleNotFound):
			utils.ServerError(c, "Default user role is not configured")
		default:
			utils.ServerError(c, "Unable to register user")
		}
		return
	}

	utils.Created(c, "User registered successfully", result)
}

// Login godoc
// @Summary      Log in
// @Description  Authenticates a user and returns an auth token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginRequest  true  "Login payload"
// @Success      200      {object}  utils.APIResponse{data=dto.AuthResponse}
// @Failure      401      {object}  utils.APIResponse
// @Failure      422      {object}  utils.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Unauthenticated(c, "Invalid email or password")
			return
		}
		utils.ServerError(c, "Unable to login")
		return
	}

	utils.Success(c, http.StatusOK, "Login successful", result)
}

// Me godoc
// @Summary      Current user
// @Description  Returns the authenticated user's profile.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=dto.UserResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authentication required")
		return
	}

	result, err := h.authService.Me(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.NotFound(c, "User not found")
			return
		}
		utils.ServerError(c, "Unable to fetch profile")
		return
	}

	utils.Success(c, http.StatusOK, "Current user fetched", result)
}

// Logout godoc
// @Summary      Log out
// @Description  Revokes the current bearer token so it can no longer be used.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	expiresAt, _ := middleware.TokenExpiresAt(c)
	h.authService.RevokeToken(middleware.Token(c), expiresAt)

	utils.Success(c, http.StatusOK, "Logout successful", gin.H{"revoked": true})
}

// RegistrationData godoc
// @Summary      Registration reference data
// @Description  Returns reference data for the sign-up screen, including roles.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=dto.RegistrationDataResponse}
// @Router       /auth/registration-data [get]
func (h *AuthHandler) RegistrationData(c *gin.Context) {
	data, err := h.authService.RegistrationData(c.Request.Context())
	if err != nil {
		utils.ServerError(c, "Unable to load registration data")
		return
	}

	utils.Success(c, http.StatusOK, "Registration data fetched", data)
}
