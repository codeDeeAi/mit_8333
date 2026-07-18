package handler

import (
	"net/http"

	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// CategoryHandler exposes request category endpoints.
type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List godoc
// @Summary      List request categories
// @Description  Returns the maintenance request categories used by the create form.
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=[]dto.CategoryResponse}
// @Router       /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.List(c.Request.Context())
	if err != nil {
		utils.ServerError(c, "Unable to fetch categories")
		return
	}

	utils.Success(c, http.StatusOK, "Categories fetched successfully", categories)
}
