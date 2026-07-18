package handler

import (
	"net/http"

	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// ReportHandler exposes reporting endpoints.
type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// Summary godoc
// @Summary      Reports summary
// @Description  Admin dashboard totals and counts by status and category.
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=dto.ReportSummaryResponse}
// @Router       /reports/summary [get]
func (h *ReportHandler) Summary(c *gin.Context) {
	summary, err := h.reportService.Summary(c.Request.Context())
	if err != nil {
		utils.ServerError(c, "Unable to build report summary")
		return
	}

	utils.Success(c, http.StatusOK, "Report summary fetched successfully", summary)
}
