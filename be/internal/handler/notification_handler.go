package handler

import (
	"errors"
	"net/http"
	"strconv"

	"UMSRMS/internal/middleware"
	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// NotificationHandler exposes in-app notification endpoints.
type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// List godoc
// @Summary      List notifications
// @Description  Returns the authenticated user's in-app notifications.
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=[]dto.NotificationResponse}
// @Router       /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authentication required")
		return
	}

	notifications, err := h.notificationService.List(c.Request.Context(), userID)
	if err != nil {
		utils.ServerError(c, "Unable to fetch notifications")
		return
	}

	utils.Success(c, http.StatusOK, "Notifications fetched successfully", notifications)
}

// MarkRead godoc
// @Summary      Mark a notification read
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Notification ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /notifications/{id}/read [put]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authentication required")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		utils.BadRequest(c, "Invalid notification id", gin.H{"id": c.Param("id")})
		return
	}

	if err := h.notificationService.MarkRead(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			utils.NotFound(c, "Notification not found")
			return
		}
		utils.ServerError(c, "Unable to update notification")
		return
	}

	utils.Success(c, http.StatusOK, "Notification marked as read", gin.H{"read": true})
}

// MarkAllRead godoc
// @Summary      Mark all notifications read
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse
// @Router       /notifications/read-all [put]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		utils.Unauthenticated(c, "Authentication required")
		return
	}

	if err := h.notificationService.MarkAllRead(c.Request.Context(), userID); err != nil {
		utils.ServerError(c, "Unable to update notifications")
		return
	}

	utils.Success(c, http.StatusOK, "All notifications marked as read", gin.H{"read": true})
}
