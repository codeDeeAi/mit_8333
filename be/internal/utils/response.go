package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful API response.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created API response.
func Created(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Resource created successfully"
	}
	Success(c, http.StatusCreated, message, data)
}

// Error sends a failed API response.
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest sends a 400 bad request response.
func BadRequest(c *gin.Context, message string, err interface{}) {
	if message == "" {
		message = "Bad request"
	}
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthenticated sends a 401 unauthenticated response.
func Unauthenticated(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Error(c, http.StatusUnauthorized, message, gin.H{"code": "UNAUTHENTICATED"})
}

// Unauthorized sends a 403 unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "You are not allowed to perform this action"
	}
	Error(c, http.StatusForbidden, message, gin.H{"code": "UNAUTHORIZED"})
}

// ValidationError sends a 422 validation error response.
func ValidationError(c *gin.Context, message string, details interface{}) {
	if message == "" {
		message = "Validation failed"
	}
	Error(c, http.StatusUnprocessableEntity, message, gin.H{
		"code":    "VALIDATION_ERROR",
		"details": details,
	})
}

// NotFound sends a 404 not found response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message, gin.H{"code": "NOT_FOUND"})
}

// Conflict sends a 409 conflict response.
func Conflict(c *gin.Context, message string, err interface{}) {
	if message == "" {
		message = "Resource conflict"
	}
	Error(c, http.StatusConflict, message, err)
}

// TooManyRequests sends a 429 rate-limit response.
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = "Too many requests, please slow down"
	}
	Error(c, http.StatusTooManyRequests, message, gin.H{"code": "RATE_LIMITED"})
}

// ServerError sends a 500 internal server error response.
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Something went wrong"
	}
	Error(c, http.StatusInternalServerError, message, gin.H{"code": "INTERNAL_SERVER_ERROR"})
}
