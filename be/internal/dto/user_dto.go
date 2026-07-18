package dto

// UpdateUserRoleRequest is the payload for changing a user's role.
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student_staff maintenance_officer admin"`
}
