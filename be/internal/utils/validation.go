package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidationTagNames makes validation errors report the json field name
// (e.g. "full_name") instead of the Go struct field name ("FullName").
func RegisterValidationTagNames() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return field.Name
			}
			return name
		})
	}
}

// FormatValidationErrors turns validator errors into a field -> message map.
// It returns false when err is not a validation error (e.g. malformed JSON).
func FormatValidationErrors(err error) (map[string]string, bool) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil, false
	}

	out := make(map[string]string, len(ve))
	for _, fe := range ve {
		out[fe.Field()] = messageForError(fe)
	}
	return out, true
}

func messageForError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return "Invalid value"
	}
}
