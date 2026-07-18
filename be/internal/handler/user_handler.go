package handler

import (
	"errors"
	"net/http"
	"strconv"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/middleware"
	"UMSRMS/internal/models"
	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler exposes admin user management endpoints.
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List godoc
// @Summary      List users
// @Description  Admins list all users, optionally filtered by role.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        role  query  string  false  "Filter by role"  Enums(student_staff, maintenance_officer, admin)
// @Success      200  {object}  utils.APIResponse{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context(), c.Query("role"))
	if err != nil {
		utils.ServerError(c, "Unable to fetch users")
		return
	}

	utils.Success(c, http.StatusOK, "Users fetched successfully", users)
}

// ListOfficers godoc
// @Summary      List maintenance officers
// @Description  Returns users with the maintenance officer role (for assignment).
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=[]dto.UserResponse}
// @Router       /users/officers [get]
func (h *UserHandler) ListOfficers(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context(), models.RoleMaintenanceOfficer)
	if err != nil {
		utils.ServerError(c, "Unable to fetch officers")
		return
	}

	utils.Success(c, http.StatusOK, "Officers fetched successfully", users)
}

// UpdateRole godoc
// @Summary      Update a user's role
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int                        true  "User ID"
// @Param        payload  body  dto.UpdateUserRoleRequest  true  "New role"
// @Success      200  {object}  utils.APIResponse{data=dto.UserResponse}
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Router       /users/{id}/role [put]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	actorID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if details, ok := utils.FormatValidationErrors(err); ok {
			utils.ValidationError(c, "Validation failed", details)
			return
		}
		utils.BadRequest(c, "Invalid role payload", err.Error())
		return
	}

	result, err := h.userService.UpdateRole(c.Request.Context(), actorID, userID, req.Role)
	if err != nil {
		handleUserError(c, err, "Unable to update user role")
		return
	}

	utils.Success(c, http.StatusOK, "User role updated successfully", result)
}

// Delete godoc
// @Summary      Delete a user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "User ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	actorID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	if err := h.userService.Delete(c.Request.Context(), actorID, userID); err != nil {
		handleUserError(c, err, "Unable to delete user")
		return
	}

	utils.Success(c, http.StatusOK, "User deleted successfully", gin.H{"deleted": true})
}

func parseUserID(c *gin.Context) (int64, bool) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userID <= 0 {
		utils.BadRequest(c, "Invalid user id", gin.H{"id": c.Param("id")})
		return 0, false
	}
	return userID, true
}

func handleUserError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		utils.NotFound(c, "User not found")
	case errors.Is(err, service.ErrRoleNotFound):
		utils.BadRequest(c, "Invalid role", nil)
	case errors.Is(err, service.ErrCannotModifySelf):
		utils.Unauthorized(c, "You cannot modify your own account")
	default:
		utils.ServerError(c, fallbackMessage)
	}
}
