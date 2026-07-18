package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"UMSRMS/internal/models"

	squirrel "github.com/Masterminds/squirrel"
)

// ServiceRequestListFilter defines supported request list filters and scoping.
type ServiceRequestListFilter struct {
	UserID        int64
	Role          string
	Search        string
	Status        string
	CategoryID    int64
	Priority      string
	Page          int
	PageSize      int
	SortColumn    string
	SortDirection string
}

// ServiceRequestWithMeta is a request enriched with joined display data.
type ServiceRequestWithMeta struct {
	models.ServiceRequest
	CategoryName        string
	CreatedByName       string
	AssignedOfficerID   *int64
	AssignedOfficerName *string
}

var serviceRequestMetaColumns = []string{
	"sr.id",
	"sr.title",
	"sr.description",
	"sr.category_id",
	"sr.created_by",
	"sr.location",
	"sr.priority",
	"sr.status",
	"sr.evidence_path",
	"sr.created_at",
	"sr.updated_at",
	"rc.name AS category_name",
	"cu.full_name AS created_by_name",
	"la.officer_id AS assigned_officer_id",
	"ou.full_name AS assigned_officer_name",
}

// ServiceRequestRepository handles service request persistence.
type ServiceRequestRepository struct {
	db *sql.DB
}

func NewServiceRequestRepository(db *sql.DB) *ServiceRequestRepository {
	return &ServiceRequestRepository{db: db}
}

func (r *ServiceRequestRepository) Create(ctx context.Context, request *models.ServiceRequest) (*models.ServiceRequest, error) {
	query := SQL.Insert("service_requests").
		Columns("title", "description", "category_id", "created_by", "location", "priority", "status", "evidence_path").
		Values(request.Title, request.Description, request.CategoryID, request.CreatedBy, request.Location, request.Priority, request.Status, request.EvidencePath).
		Suffix("RETURNING id, title, description, category_id, created_by, location, priority, status, evidence_path, created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	created := &models.ServiceRequest{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.CategoryID,
		&created.CreatedBy,
		&created.Location,
		&created.Priority,
		&created.Status,
		&created.EvidencePath,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *ServiceRequestRepository) GetByID(ctx context.Context, id int64) (*models.ServiceRequest, error) {
	query := SQL.Select(
		"id",
		"title",
		"description",
		"category_id",
		"created_by",
		"location",
		"priority",
		"status",
		"evidence_path",
		"created_at",
		"updated_at",
	).From("service_requests").Where("id = ?", id).Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	request := &models.ServiceRequest{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&request.ID,
		&request.Title,
		&request.Description,
		&request.CategoryID,
		&request.CreatedBy,
		&request.Location,
		&request.Priority,
		&request.Status,
		&request.EvidencePath,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *ServiceRequestRepository) List(ctx context.Context, filter ServiceRequestListFilter) ([]ServiceRequestWithMeta, int64, error) {
	baseSelect := applyServiceRequestJoins(SQL.Select(serviceRequestMetaColumns...).From("service_requests sr"))
	countQuery := SQL.Select("COUNT(sr.id)").From("service_requests sr")

	baseSelect, countQuery = applyServiceRequestScope(baseSelect, countQuery, filter)
	baseSelect, countQuery = applyServiceRequestFilters(baseSelect, countQuery, filter)

	baseSelect = baseSelect.OrderBy(fmt.Sprintf("sr.%s %s", filter.SortColumn, filter.SortDirection))
	baseSelect = baseSelect.Limit(uint64(filter.PageSize)).Offset(uint64((filter.Page - 1) * filter.PageSize))

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL, listArgs, err := baseSelect.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	requests := make([]ServiceRequestWithMeta, 0)
	for rows.Next() {
		request, err := scanServiceRequestMeta(rows)
		if err != nil {
			return nil, 0, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetDetailByID returns a request enriched with category, creator and current
// officer names.
func (r *ServiceRequestRepository) GetDetailByID(ctx context.Context, id int64) (*ServiceRequestWithMeta, error) {
	query := applyServiceRequestJoins(SQL.Select(serviceRequestMetaColumns...).From("service_requests sr")).
		Where("sr.id = ?", id).
		Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	item, err := scanServiceRequestMeta(r.db.QueryRowContext(ctx, sqlQuery, args...))
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Assign records an officer assignment and moves a pending request to assigned.
func (r *ServiceRequestRepository) Assign(ctx context.Context, requestID, officerID, assignedBy int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statusQuery, statusArgs, err := SQL.Select("status").From("service_requests").Where("id = ?", requestID).Suffix("FOR UPDATE").ToSql()
	if err != nil {
		return err
	}
	var status string
	if err := tx.QueryRowContext(ctx, statusQuery, statusArgs...).Scan(&status); err != nil {
		return err
	}

	insertQuery, insertArgs, err := SQL.Insert("assignments").
		Columns("request_id", "officer_id", "assigned_by").
		Values(requestID, officerID, assignedBy).
		ToSql()
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
		return err
	}

	if status == models.StatusPending {
		updateQuery, updateArgs, err := SQL.Update("service_requests").
			Set("status", models.StatusAssigned).
			Set("updated_at", squirrel.Expr("NOW()")).
			Where("id = ?", requestID).
			ToSql()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, updateQuery, updateArgs...); err != nil {
			return err
		}

		logQuery, logArgs, err := SQL.Insert("status_logs").
			Columns("request_id", "changed_by", "old_status", "new_status", "note").
			Values(requestID, assignedBy, models.StatusPending, models.StatusAssigned, nil).
			ToSql()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, logQuery, logArgs...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// OfficerExists reports whether the id belongs to a maintenance officer.
func (r *ServiceRequestRepository) OfficerExists(ctx context.Context, officerID int64) (bool, error) {
	query := SQL.Select("1").
		From("users u").
		Join("roles r ON r.id = u.role_id").
		Where("u.id = ?", officerID).
		Where("r.name = ?", models.RoleMaintenanceOfficer).
		Limit(1)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	if err := r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func applyServiceRequestJoins(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	return query.
		LeftJoin("request_categories rc ON rc.id = sr.category_id").
		LeftJoin("users cu ON cu.id = sr.created_by").
		JoinClause("LEFT JOIN LATERAL (SELECT officer_id FROM assignments WHERE request_id = sr.id ORDER BY assigned_at DESC LIMIT 1) la ON true").
		LeftJoin("users ou ON ou.id = la.officer_id")
}

func scanServiceRequestMeta(scanner interface{ Scan(...any) error }) (ServiceRequestWithMeta, error) {
	var item ServiceRequestWithMeta
	var categoryName, createdByName, officerName sql.NullString
	var officerID sql.NullInt64

	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.CategoryID,
		&item.CreatedBy,
		&item.Location,
		&item.Priority,
		&item.Status,
		&item.EvidencePath,
		&item.CreatedAt,
		&item.UpdatedAt,
		&categoryName,
		&createdByName,
		&officerID,
		&officerName,
	); err != nil {
		return item, err
	}

	item.CategoryName = categoryName.String
	item.CreatedByName = createdByName.String
	if officerID.Valid {
		id := officerID.Int64
		item.AssignedOfficerID = &id
	}
	if officerName.Valid {
		name := officerName.String
		item.AssignedOfficerName = &name
	}

	return item, nil
}

func (r *ServiceRequestRepository) IsAssignedToOfficer(ctx context.Context, requestID, officerID int64) (bool, error) {
	query := SQL.Select("1").From("assignments").Where("request_id = ?", requestID).Where("officer_id = ?", officerID).Limit(1)
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ServiceRequestRepository) GetStatusLogs(ctx context.Context, requestID int64) ([]models.StatusLog, error) {
	query := SQL.Select("id", "request_id", "changed_by", "old_status", "new_status", "note", "created_at").
		From("status_logs").
		Where("request_id = ?", requestID).
		OrderBy("created_at DESC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]models.StatusLog, 0)
	for rows.Next() {
		var log models.StatusLog
		if err := rows.Scan(&log.ID, &log.RequestID, &log.ChangedBy, &log.OldStatus, &log.NewStatus, &log.Note, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *ServiceRequestRepository) UpdateStatusWithLog(ctx context.Context, requestID, changedBy int64, newStatus string, note *string) (*models.ServiceRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	currentStatusQuery, currentStatusArgs, err := SQL.Select("status").From("service_requests").Where("id = ?", requestID).Suffix("FOR UPDATE").ToSql()
	if err != nil {
		return nil, err
	}

	var oldStatus string
	if err := tx.QueryRowContext(ctx, currentStatusQuery, currentStatusArgs...).Scan(&oldStatus); err != nil {
		return nil, err
	}

	updateQuery, updateArgs, err := SQL.Update("service_requests").
		Set("status", newStatus).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where("id = ?", requestID).
		Suffix("RETURNING id, title, description, category_id, created_by, location, priority, status, evidence_path, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	updated := &models.ServiceRequest{}
	if err := tx.QueryRowContext(ctx, updateQuery, updateArgs...).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.CategoryID,
		&updated.CreatedBy,
		&updated.Location,
		&updated.Priority,
		&updated.Status,
		&updated.EvidencePath,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	); err != nil {
		return nil, err
	}

	insertLogQuery, insertLogArgs, err := SQL.Insert("status_logs").
		Columns("request_id", "changed_by", "old_status", "new_status", "note").
		Values(requestID, changedBy, oldStatus, newStatus, note).
		ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, insertLogQuery, insertLogArgs...); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *ServiceRequestRepository) UpdateEvidencePath(ctx context.Context, requestID int64, evidencePath *string) (*models.ServiceRequest, error) {
	query := SQL.Update("service_requests").
		Set("evidence_path", evidencePath).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where("id = ?", requestID).
		Suffix("RETURNING id, title, description, category_id, created_by, location, priority, status, evidence_path, created_at, updated_at")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	updated := &models.ServiceRequest{}
	err = r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.CategoryID,
		&updated.CreatedBy,
		&updated.Location,
		&updated.Priority,
		&updated.Status,
		&updated.EvidencePath,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *ServiceRequestRepository) Delete(ctx context.Context, requestID int64) error {
	query := SQL.Delete("service_requests").Where("id = ?", requestID)
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func applyServiceRequestScope(selectQuery, countQuery squirrel.SelectBuilder, filter ServiceRequestListFilter) (squirrel.SelectBuilder, squirrel.SelectBuilder) {
	switch filter.Role {
	case models.RoleStudentStaff:
		selectQuery = selectQuery.Where("sr.created_by = ?", filter.UserID)
		countQuery = countQuery.Where("sr.created_by = ?", filter.UserID)
	case models.RoleMaintenanceOfficer:
		predicate := squirrel.Expr("EXISTS (SELECT 1 FROM assignments a WHERE a.request_id = sr.id AND a.officer_id = ?)", filter.UserID)
		selectQuery = selectQuery.Where(predicate)
		countQuery = countQuery.Where(predicate)
	}

	return selectQuery, countQuery
}

func applyServiceRequestFilters(selectQuery, countQuery squirrel.SelectBuilder, filter ServiceRequestListFilter) (squirrel.SelectBuilder, squirrel.SelectBuilder) {
	if search := strings.TrimSpace(filter.Search); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		predicate := squirrel.Expr("(LOWER(sr.title) LIKE ? OR LOWER(sr.description) LIKE ? OR LOWER(sr.location) LIKE ?)", like, like, like)
		selectQuery = selectQuery.Where(predicate)
		countQuery = countQuery.Where(predicate)
	}

	if filter.Status != "" {
		selectQuery = selectQuery.Where("sr.status = ?", filter.Status)
		countQuery = countQuery.Where("sr.status = ?", filter.Status)
	}
	if filter.CategoryID > 0 {
		selectQuery = selectQuery.Where("sr.category_id = ?", filter.CategoryID)
		countQuery = countQuery.Where("sr.category_id = ?", filter.CategoryID)
	}
	if filter.Priority != "" {
		selectQuery = selectQuery.Where("sr.priority = ?", filter.Priority)
		countQuery = countQuery.Where("sr.priority = ?", filter.Priority)
	}

	return selectQuery, countQuery
}
