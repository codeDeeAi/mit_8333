package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"UMSRMS/internal/config"
	"UMSRMS/internal/dto"
	"UMSRMS/internal/models"
	"UMSRMS/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrRequestNotFound         = errors.New("service request not found")
	ErrRequestForbidden        = errors.New("you are not allowed to access this request")
	ErrUnsupportedRequestRole  = errors.New("unsupported request access role")
	ErrEvidenceRequired        = errors.New("evidence file is required")
	ErrEvidenceTooLarge        = errors.New("evidence file exceeds size limit")
	ErrEvidenceTypeUnsupported = errors.New("unsupported evidence file type")
	ErrInvalidRequestInput     = errors.New("invalid request input")
	ErrOfficerNotFound         = errors.New("maintenance officer not found")
)

type RequestActor struct {
	UserID int64
	Role   string
}

type ServiceRequestService struct {
	repo        *repository.ServiceRequestRepository
	notifier    *repository.NotificationRepository
	uploadDir   string
	maxUploadMB int
}

func NewServiceRequestService(repo *repository.ServiceRequestRepository, notifier *repository.NotificationRepository, cfg *config.EnvConfig) *ServiceRequestService {
	return &ServiceRequestService{
		repo:        repo,
		notifier:    notifier,
		uploadDir:   cfg.UploadDir,
		maxUploadMB: cfg.MaxUploadMB,
	}
}

func (s *ServiceRequestService) Create(ctx context.Context, actor RequestActor, req dto.CreateServiceRequestRequest, evidence *multipart.FileHeader) (*dto.ServiceRequestResponse, error) {
	if actor.Role != models.RoleStudentStaff {
		return nil, ErrRequestForbidden
	}

	priority := req.Priority
	if priority == "" {
		priority = models.PriorityMedium
	}

	var evidencePath *string
	if evidence != nil {
		savedPath, err := s.saveEvidenceFile(evidence)
		if err != nil {
			return nil, err
		}
		evidencePath = &savedPath
	}

	created, err := s.repo.Create(ctx, &models.ServiceRequest{
		Title:        req.Title,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		CreatedBy:    actor.UserID,
		Location:     req.Location,
		Priority:     priority,
		Status:       models.StatusPending,
		EvidencePath: evidencePath,
	})
	if err != nil {
		if evidencePath != nil {
			_ = os.Remove(*evidencePath)
		}
		return nil, normalizeRequestRepositoryError(err)
	}

	_ = s.notifier.CreateForRole(ctx, models.RoleAdmin, "New service request: "+created.Title)

	response := toServiceRequestResponse(created)
	return &response, nil
}

func (s *ServiceRequestService) List(ctx context.Context, actor RequestActor, query dto.ListServiceRequestsQuery) (*dto.ServiceRequestListResponse, error) {
	if !isSupportedRequestRole(actor.Role) {
		return nil, ErrUnsupportedRequestRole
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sortColumn, sortDirection := parseRequestSort(query.Sort)

	items, total, err := s.repo.List(ctx, repository.ServiceRequestListFilter{
		UserID:        actor.UserID,
		Role:          actor.Role,
		Search:        query.Q,
		Status:        query.Status,
		CategoryID:    query.CategoryID,
		Priority:      query.Priority,
		Page:          page,
		PageSize:      pageSize,
		SortColumn:    sortColumn,
		SortDirection: sortDirection,
	})
	if err != nil {
		return nil, err
	}

	responseItems := make([]dto.ServiceRequestResponse, 0, len(items))
	for i := range items {
		responseItems = append(responseItems, toServiceRequestResponseMeta(&items[i]))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &dto.ServiceRequestListResponse{
		Items: responseItems,
		Pagination: dto.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ServiceRequestService) GetByID(ctx context.Context, actor RequestActor, requestID int64) (*dto.ServiceRequestDetailResponse, error) {
	request, err := s.repo.GetDetailByID(ctx, requestID)
	if err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	allowed, err := s.canAccessRequest(ctx, actor, &request.ServiceRequest)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrRequestForbidden
	}

	logs, err := s.repo.GetStatusLogs(ctx, requestID)
	if err != nil {
		return nil, err
	}

	responseLogs := make([]dto.StatusLogResponse, 0, len(logs))
	for _, log := range logs {
		responseLogs = append(responseLogs, dto.StatusLogResponse{
			ID:        log.ID,
			RequestID: log.RequestID,
			ChangedBy: log.ChangedBy,
			OldStatus: log.OldStatus,
			NewStatus: log.NewStatus,
			Note:      log.Note,
			CreatedAt: log.CreatedAt,
		})
	}

	return &dto.ServiceRequestDetailResponse{
		Request:    toServiceRequestResponseMeta(request),
		StatusLogs: responseLogs,
	}, nil
}

func (s *ServiceRequestService) UpdateStatus(ctx context.Context, actor RequestActor, requestID int64, req dto.UpdateServiceRequestStatusRequest) (*dto.ServiceRequestResponse, error) {
	request, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	switch actor.Role {
	case models.RoleAdmin:
	case models.RoleMaintenanceOfficer:
		assigned, err := s.repo.IsAssignedToOfficer(ctx, request.ID, actor.UserID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, ErrRequestForbidden
		}
	default:
		return nil, ErrRequestForbidden
	}

	updated, err := s.repo.UpdateStatusWithLog(ctx, requestID, actor.UserID, req.Status, req.Note)
	if err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	if updated.CreatedBy != actor.UserID {
		_ = s.notifier.Create(ctx, updated.CreatedBy, "Your request \""+updated.Title+"\" is now "+humanStatus(updated.Status))
	}

	response := toServiceRequestResponse(updated)
	return &response, nil
}

func (s *ServiceRequestService) Assign(ctx context.Context, actor RequestActor, requestID, officerID int64) (*dto.ServiceRequestResponse, error) {
	if actor.Role != models.RoleAdmin {
		return nil, ErrRequestForbidden
	}

	if _, err := s.repo.GetByID(ctx, requestID); err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	isOfficer, err := s.repo.OfficerExists(ctx, officerID)
	if err != nil {
		return nil, err
	}
	if !isOfficer {
		return nil, ErrOfficerNotFound
	}

	if err := s.repo.Assign(ctx, requestID, officerID, actor.UserID); err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	detail, err := s.repo.GetDetailByID(ctx, requestID)
	if err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	_ = s.notifier.Create(ctx, officerID, "You have been assigned: \""+detail.Title+"\"")
	if detail.CreatedBy != officerID {
		_ = s.notifier.Create(ctx, detail.CreatedBy, "Your request \""+detail.Title+"\" was assigned to a maintenance officer")
	}

	response := toServiceRequestResponseMeta(detail)
	return &response, nil
}

func (s *ServiceRequestService) Delete(ctx context.Context, actor RequestActor, requestID int64) error {
	if actor.Role != models.RoleAdmin {
		return ErrRequestForbidden
	}

	request, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return normalizeRequestRepositoryError(err)
	}

	if err := s.repo.Delete(ctx, requestID); err != nil {
		return normalizeRequestRepositoryError(err)
	}

	if request.EvidencePath != nil {
		_ = os.Remove(*request.EvidencePath)
	}

	return nil
}

func (s *ServiceRequestService) UploadEvidence(ctx context.Context, actor RequestActor, requestID int64, evidence *multipart.FileHeader) (*dto.ServiceRequestResponse, error) {
	if evidence == nil {
		return nil, ErrEvidenceRequired
	}

	request, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return nil, normalizeRequestRepositoryError(err)
	}

	if request.CreatedBy != actor.UserID {
		return nil, ErrRequestForbidden
	}

	newPath, err := s.saveEvidenceFile(evidence)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateEvidencePath(ctx, requestID, &newPath)
	if err != nil {
		_ = os.Remove(newPath)
		return nil, normalizeRequestRepositoryError(err)
	}

	if request.EvidencePath != nil {
		_ = os.Remove(*request.EvidencePath)
	}

	response := toServiceRequestResponse(updated)
	return &response, nil
}

func (s *ServiceRequestService) canAccessRequest(ctx context.Context, actor RequestActor, request *models.ServiceRequest) (bool, error) {
	switch actor.Role {
	case models.RoleAdmin:
		return true, nil
	case models.RoleStudentStaff:
		return request.CreatedBy == actor.UserID, nil
	case models.RoleMaintenanceOfficer:
		return s.repo.IsAssignedToOfficer(ctx, request.ID, actor.UserID)
	default:
		return false, ErrUnsupportedRequestRole
	}
}

func (s *ServiceRequestService) saveEvidenceFile(fileHeader *multipart.FileHeader) (string, error) {
	maxBytes := int64(s.maxUploadMB) * 1024 * 1024
	if fileHeader.Size > maxBytes {
		return "", ErrEvidenceTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:bytesRead])
	if !isAllowedEvidenceType(contentType) {
		return "", ErrEvidenceTypeUnsupported
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return "", err
		}
	} else {
		return "", errors.New("unable to reset uploaded file stream")
	}

	directory := filepath.Join(s.uploadDir, "requests")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", err
	}

	filename, err := randomFilename(filepath.Ext(fileHeader.Filename))
	if err != nil {
		return "", err
	}

	path := filepath.Join(directory, filename)
	destination, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return "", err
	}

	return path, nil
}

func randomFilename(ext string) (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer) + strings.ToLower(ext), nil
}

func isAllowedEvidenceType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/webp", "application/pdf":
		return true
	default:
		return false
	}
}

func parseRequestSort(sort string) (string, string) {
	column := "created_at"
	direction := "DESC"

	if sort == "" {
		return column, direction
	}

	parts := strings.Split(sort, ".")
	if len(parts) > 0 {
		switch parts[0] {
		case "created_at", "updated_at", "priority", "status", "title":
			column = parts[0]
		}
	}
	if len(parts) > 1 && strings.EqualFold(parts[1], "asc") {
		direction = "ASC"
	}

	return column, direction
}

func isSupportedRequestRole(role string) bool {
	switch role {
	case models.RoleStudentStaff, models.RoleMaintenanceOfficer, models.RoleAdmin:
		return true
	default:
		return false
	}
}

func toServiceRequestResponse(request *models.ServiceRequest) dto.ServiceRequestResponse {
	return dto.ServiceRequestResponse{
		ID:           request.ID,
		Title:        request.Title,
		Description:  request.Description,
		CategoryID:   request.CategoryID,
		CreatedBy:    request.CreatedBy,
		Location:     request.Location,
		Priority:     request.Priority,
		Status:       request.Status,
		EvidencePath: request.EvidencePath,
		CreatedAt:    request.CreatedAt,
		UpdatedAt:    request.UpdatedAt,
	}
}

func humanStatus(status string) string {
	return strings.ReplaceAll(status, "_", " ")
}

func toServiceRequestResponseMeta(m *repository.ServiceRequestWithMeta) dto.ServiceRequestResponse {
	response := toServiceRequestResponse(&m.ServiceRequest)
	response.CategoryName = m.CategoryName
	response.CreatedByName = m.CreatedByName
	response.AssignedOfficerID = m.AssignedOfficerID
	response.AssignedOfficerName = m.AssignedOfficerName
	return response
}

func normalizeRequestRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRequestNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503", "23514", "22P02":
			return ErrInvalidRequestInput
		}
	}

	return err
}
