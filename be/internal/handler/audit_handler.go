package handler

import (
	"net/http"

	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuditHandler exposes the audit trail endpoint.
type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List godoc
// @Summary      List audit log
// @Description  Admins view the activity/audit trail of key actions.
// @Tags         audit
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=[]dto.AuditLogResponse}
// @Router       /audit-logs [get]
func (h *AuditHandler) List(c *gin.Context) {
	logs, err := h.auditService.List(c.Request.Context())
	if err != nil {
		utils.ServerError(c, "Unable to fetch audit logs")
		return
	}

	utils.Success(c, http.StatusOK, "Audit logs fetched successfully", logs)
}
