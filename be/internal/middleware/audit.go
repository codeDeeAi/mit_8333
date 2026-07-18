package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// AuditRecorder persists audit entries. Implemented by the audit repository.
type AuditRecorder interface {
	Record(ctx context.Context, userID int64, action, entity, entityID string) error
}

type auditRule struct {
	action string
	entity string
}

// auditRules maps "METHOD route-pattern" to a semantic audit action.
var auditRules = map[string]auditRule{
	"POST /api/v1/requests":              {"CREATE_REQUEST", "service_request"},
	"POST /api/v1/requests/:id/assign":   {"ASSIGN_REQUEST", "service_request"},
	"PUT /api/v1/requests/:id/status":    {"UPDATE_STATUS", "service_request"},
	"DELETE /api/v1/requests/:id":        {"DELETE_REQUEST", "service_request"},
	"POST /api/v1/requests/:id/evidence": {"UPLOAD_EVIDENCE", "service_request"},
	"PUT /api/v1/users/:id/role":         {"UPDATE_USER_ROLE", "user"},
	"DELETE /api/v1/users/:id":           {"DELETE_USER", "user"},
}

// AuditLog records successful mutating requests to the audit trail. It runs
// after the handler and only logs 2xx/3xx responses on mapped routes.
func AuditLog(recorder AuditRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() >= 400 {
			return
		}

		rule, ok := auditRules[c.Request.Method+" "+c.FullPath()]
		if !ok {
			return
		}

		userID, _ := UserID(c)
		// Detached context: the request context may be done once the response is written.
		_ = recorder.Record(context.Background(), userID, rule.action, rule.entity, c.Param("id"))
	}
}
