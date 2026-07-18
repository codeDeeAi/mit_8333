package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"UMSRMS/internal/config"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newRouter(cfg *config.EnvConfig) *gin.Engine {
	r := gin.New()
	r.Use(CORS(cfg))
	r.Use(RateLimit(cfg))
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	return r
}

func TestCORSAllowAllOrigins(t *testing.T) {
	cfg := &config.EnvConfig{CORSAllowedOrigins: []string{"*"}}
	r := newRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://any-domain.example")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected Access-Control-Allow-Origin '*', got %q", got)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCORSPreflightShortCircuits(t *testing.T) {
	cfg := &config.EnvConfig{CORSAllowedOrigins: []string{"*"}}
	r := newRouter(cfg)

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://any-domain.example")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("expected Access-Control-Allow-Methods to be set")
	}
}

func TestCORSAllowlistRejectsUnknownOrigin(t *testing.T) {
	cfg := &config.EnvConfig{CORSAllowedOrigins: []string{"http://localhost:5173"}}
	r := newRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allow-origin for unlisted origin, got %q", got)
	}
}

func TestRateLimitReturns429WhenExceeded(t *testing.T) {
	// rps 1, burst 1 -> first request passes, immediate second is throttled.
	cfg := &config.EnvConfig{
		CORSAllowedOrigins: []string{"*"},
		RateLimitEnabled:   true,
		RateLimitRPS:       1,
		RateLimitBurst:     1,
	}
	r := newRouter(cfg)

	do := func() int {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.RemoteAddr = "203.0.113.7:5555"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	if code := do(); code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", code)
	}
	if code := do(); code != http.StatusTooManyRequests {
		t.Fatalf("second request: expected 429, got %d", code)
	}
}

func TestRateLimitPerMinuteAllowsBurstThenBlocks(t *testing.T) {
	r := gin.New()
	r.Use(RateLimitPerMinute(5))
	r.GET("/auth/login", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	do := func() int {
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		req.RemoteAddr = "198.51.100.4:5555"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	// First 5 requests consume the burst and succeed.
	for i := 1; i <= 5; i++ {
		if code := do(); code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, code)
		}
	}
	// The 6th within the same minute is throttled.
	if code := do(); code != http.StatusTooManyRequests {
		t.Fatalf("6th request: expected 429, got %d", code)
	}
}

func TestRequireAuthValidTokenSetsContext(t *testing.T) {
	cfg := &config.EnvConfig{AppName: "test", JWTSecret: "test-secret", JWTExpireHours: 24}
	jwtManager, err := utils.NewJWTManager(cfg)
	if err != nil {
		t.Fatalf("jwt manager: %v", err)
	}
	banList := utils.NewTokenBanList()
	token, _, err := jwtManager.GenerateToken("42", "user@miva.edu", "admin")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := gin.New()
	r.Use(RequireAuth(jwtManager, banList))
	r.GET("/auth/logout", func(c *gin.Context) {
		id, ok := UserID(c)
		if !ok || id != 42 {
			t.Errorf("expected user id 42 in context, got %d (ok=%v)", id, ok)
		}
		if Token(c) != token {
			t.Error("expected raw token in context")
		}
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequireAuthRejectsMissingAndBannedTokens(t *testing.T) {
	cfg := &config.EnvConfig{AppName: "test", JWTSecret: "test-secret", JWTExpireHours: 24}
	jwtManager, _ := utils.NewJWTManager(cfg)
	banList := utils.NewTokenBanList()
	token, expiresAt, _ := jwtManager.GenerateToken("42", "user@miva.edu", "admin")

	r := gin.New()
	r.Use(RequireAuth(jwtManager, banList))
	r.GET("/protected", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	call := func(header string) int {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	if code := call(""); code != http.StatusUnauthorized {
		t.Fatalf("missing header: expected 401, got %d", code)
	}
	if code := call("Bearer not-a-real-token"); code != http.StatusUnauthorized {
		t.Fatalf("invalid token: expected 401, got %d", code)
	}
	// A banned (revoked) token must be rejected.
	banList.Ban(token, expiresAt)
	if code := call("Bearer " + token); code != http.StatusUnauthorized {
		t.Fatalf("banned token: expected 401, got %d", code)
	}
}

func TestRateLimitDisabledIsNoOp(t *testing.T) {
	cfg := &config.EnvConfig{
		CORSAllowedOrigins: []string{"*"},
		RateLimitEnabled:   false,
	}
	r := newRouter(cfg)

	for i := range 5 {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.RemoteAddr = "203.0.113.9:5555"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200 with limiter disabled, got %d", i, w.Code)
		}
	}
}
