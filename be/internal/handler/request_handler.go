package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"UMSRMS/internal/dto"
	"UMSRMS/internal/middleware"
	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

// ServiceRequestHandler exposes request CRUD endpoints.
type ServiceRequestHandler struct {
	requestService *service.ServiceRequestService
}

func NewServiceRequestHandler(requestService *service.ServiceRequestService) *ServiceRequestHandler {
	return &ServiceRequestHandler{requestService: requestService}
}

// Create godoc
// @Summary      Create a service request
// @Description  Student/staff submit a maintenance request, optionally with an evidence file.
// @Tags         requests
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        title        formData  string  true   "Title"
// @Param        description  formData  string  true   "Description"
// @Param        category_id  formData  int     true   "Category ID"
// @Param        location     formData  string  true   "Location"
// @Param        priority     formData  string  false  "Priority"  Enums(low, medium, high)
// @Param        evidence     formData  file    false  "Evidence file (image or PDF)"
// @Success      201  {object}  utils.APIResponse{data=dto.ServiceRequestResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Router       /requests [post]
func (h *ServiceRequestHandler) Create(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	var req dto.CreateServiceRequestRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationError(c, "Invalid request payload", err.Error())
		return
	}

	var evidenceHeader *multipart.FileHeader
	fileHeader, err := c.FormFile("evidence")
	if err == nil {
		evidenceHeader = fileHeader
	} else if !errors.Is(err, http.ErrMissingFile) {
		utils.BadRequest(c, "Invalid evidence upload", err.Error())
		return
	}

	result, err := h.requestService.Create(c.Request.Context(), actor, req, evidenceHeader)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to create request")
		return
	}

	utils.Created(c, "Service request created successfully", result)
}

// List godoc
// @Summary      List service requests
// @Description  Returns requests scoped to the caller's role, with search, filter and pagination.
// @Tags         requests
// @Produce      json
// @Security     BearerAuth
// @Param        q            query  string  false  "Search text"
// @Param        status       query  string  false  "Status"    Enums(pending, assigned, in_progress, completed, rejected)
// @Param        category_id  query  int     false  "Category ID"
// @Param        priority     query  string  false  "Priority"  Enums(low, medium, high)
// @Param        page         query  int     false  "Page number"
// @Param        page_size    query  int     false  "Page size"
// @Param        sort         query  string  false  "Sort, e.g. created_at.desc"
// @Success      200  {object}  utils.APIResponse{data=dto.ServiceRequestListResponse}
// @Router       /requests [get]
func (h *ServiceRequestHandler) List(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	var query dto.ListServiceRequestsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationError(c, "Invalid request query", err.Error())
		return
	}

	result, err := h.requestService.List(c.Request.Context(), actor, query)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to fetch requests")
		return
	}

	utils.Success(c, http.StatusOK, "Service requests fetched successfully", result)
}

// GetByID godoc
// @Summary      Get a service request
// @Description  Returns a request with its status history (access scoped by role).
// @Tags         requests
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Request ID"
// @Success      200  {object}  utils.APIResponse{data=dto.ServiceRequestDetailResponse}
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /requests/{id} [get]
func (h *ServiceRequestHandler) GetByID(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	requestID, ok := parseRequestID(c)
	if !ok {
		return
	}

	result, err := h.requestService.GetByID(c.Request.Context(), actor, requestID)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to fetch request")
		return
	}

	utils.Success(c, http.StatusOK, "Service request fetched successfully", result)
}

// UpdateStatus godoc
// @Summary      Update request status
// @Description  Maintenance officers (assigned) or admins move a request through its lifecycle.
// @Tags         requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int                                       true  "Request ID"
// @Param        payload  body  dto.UpdateServiceRequestStatusRequest  true  "Status update"
// @Success      200  {object}  utils.APIResponse{data=dto.ServiceRequestResponse}
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Router       /requests/{id}/status [put]
func (h *ServiceRequestHandler) UpdateStatus(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	requestID, ok := parseRequestID(c)
	if !ok {
		return
	}

	var req dto.UpdateServiceRequestStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid status update payload", err.Error())
		return
	}

	result, err := h.requestService.UpdateStatus(c.Request.Context(), actor, requestID, req)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to update request status")
		return
	}

	utils.Success(c, http.StatusOK, "Service request status updated successfully", result)
}

// Delete godoc
// @Summary      Delete a service request
// @Description  Admins delete a request and its evidence.
// @Tags         requests
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Request ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /requests/{id} [delete]
func (h *ServiceRequestHandler) Delete(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	requestID, ok := parseRequestID(c)
	if !ok {
		return
	}

	if err := h.requestService.Delete(c.Request.Context(), actor, requestID); err != nil {
		handleServiceRequestError(c, err, "Unable to delete request")
		return
	}

	utils.Success(c, http.StatusOK, "Service request deleted successfully", gin.H{"deleted": true})
}

// UploadEvidence godoc
// @Summary      Upload request evidence
// @Description  The request owner attaches or replaces the evidence file.
// @Tags         requests
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      int   true  "Request ID"
// @Param        evidence  formData  file  true  "Evidence file (image or PDF)"
// @Success      200  {object}  utils.APIResponse{data=dto.ServiceRequestResponse}
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Router       /requests/{id}/evidence [post]
func (h *ServiceRequestHandler) UploadEvidence(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	requestID, ok := parseRequestID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("evidence")
	if err != nil {
		utils.ValidationError(c, "Evidence file is required", err.Error())
		return
	}

	result, err := h.requestService.UploadEvidence(c.Request.Context(), actor, requestID, fileHeader)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to upload request evidence")
		return
	}

	utils.Success(c, http.StatusOK, "Evidence uploaded successfully", result)
}

// Assign godoc
// @Summary      Assign a request to an officer
// @Description  Admins assign a service request to a maintenance officer.
// @Tags         requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int                              true  "Request ID"
// @Param        payload  body  dto.AssignServiceRequestRequest  true  "Officer to assign"
// @Success      200  {object}  utils.APIResponse{data=dto.ServiceRequestResponse}
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      422  {object}  utils.APIResponse
// @Router       /requests/{id}/assign [post]
func (h *ServiceRequestHandler) Assign(c *gin.Context) {
	actor, ok := requestActorFromContext(c)
	if !ok {
		utils.Unauthenticated(c, "Authenticated user context is missing")
		return
	}

	requestID, ok := parseRequestID(c)
	if !ok {
		return
	}

	var req dto.AssignServiceRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid assignment payload", err.Error())
		return
	}

	result, err := h.requestService.Assign(c.Request.Context(), actor, requestID, req.OfficerID)
	if err != nil {
		handleServiceRequestError(c, err, "Unable to assign request")
		return
	}

	utils.Success(c, http.StatusOK, "Service request assigned successfully", result)
}

func requestActorFromContext(c *gin.Context) (service.RequestActor, bool) {
	userID, ok := middleware.UserID(c)
	if !ok {
		return service.RequestActor{}, false
	}
	role, ok := middleware.UserRole(c)
	if !ok {
		return service.RequestActor{}, false
	}

	return service.RequestActor{UserID: userID, Role: role}, true
}

func parseRequestID(c *gin.Context) (int64, bool) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestID <= 0 {
		utils.BadRequest(c, "Invalid request id", gin.H{"id": c.Param("id")})
		return 0, false
	}
	return requestID, true
}

func handleServiceRequestError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, service.ErrRequestNotFound):
		utils.NotFound(c, "Service request not found")
	case errors.Is(err, service.ErrOfficerNotFound):
		utils.NotFound(c, "Maintenance officer not found")
	case errors.Is(err, service.ErrRequestForbidden), errors.Is(err, service.ErrUnsupportedRequestRole):
		utils.Unauthorized(c, "You are not allowed to perform this action")
	case errors.Is(err, service.ErrEvidenceRequired), errors.Is(err, service.ErrInvalidRequestInput):
		utils.BadRequest(c, err.Error(), nil)
	case errors.Is(err, service.ErrEvidenceTooLarge), errors.Is(err, service.ErrEvidenceTypeUnsupported):
		utils.ValidationError(c, err.Error(), nil)
	default:
		utils.ServerError(c, fallbackMessage)
	}
}
