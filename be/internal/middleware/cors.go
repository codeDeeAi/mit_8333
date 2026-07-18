package middleware

import (
	"net/http"
	"strings"

	"UMSRMS/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that applies Cross-Origin Resource Sharing headers
// based on the configured allowed origins.
//
// When CORS_ALLOWED_ORIGINS contains "*" (the default), every origin is
// accepted and "Access-Control-Allow-Origin: *" is returned. When an explicit
// allowlist is configured, only matching origins are echoed back and
// credentials are permitted.
func CORS(cfg *config.EnvConfig) gin.HandlerFunc {
	allowAll := cfg.AllowAllOrigins()

	allowed := make(map[string]bool, len(cfg.CORSAllowedOrigins))
	for _, origin := range cfg.CORSAllowedOrigins {
		allowed[origin] = true
	}

	allowedMethods := strings.Join([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}, ", ")

	allowedHeaders := strings.Join([]string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		header := c.Writer.Header()

		switch {
		case allowAll:
			header.Set("Access-Control-Allow-Origin", "*")
		case origin != "" && allowed[origin]:
			// Reflect the specific origin so credentialed requests are allowed.
			header.Set("Access-Control-Allow-Origin", origin)
			header.Add("Vary", "Origin")
			header.Set("Access-Control-Allow-Credentials", "true")
		}

		header.Set("Access-Control-Allow-Methods", allowedMethods)
		header.Set("Access-Control-Allow-Headers", allowedHeaders)
		header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
		header.Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
