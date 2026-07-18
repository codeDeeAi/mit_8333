package middleware

import (
	"sync"
	"time"

	"UMSRMS/internal/config"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor tracks a single client's token-bucket limiter and last activity time.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimiter maintains a per-client-IP token bucket rate limiter.
type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
	ttl      time.Duration
}

func newIPRateLimiter(rps float64, burst int) *ipRateLimiter {
	l := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(rps),
		burst:    burst,
		ttl:      3 * time.Minute,
	}
	go l.cleanupLoop()
	return l
}

// limiterFor returns the rate limiter for the given IP, creating one if needed.
func (l *ipRateLimiter) limiterFor(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, ok := l.visitors[ip]
	if !ok {
		limiter := rate.NewLimiter(l.rps, l.burst)
		l.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupLoop periodically evicts visitors that have been idle beyond the TTL,
// keeping memory bounded for long-running servers.
func (l *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.ttl {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimit returns a middleware that throttles requests per client IP using a
// token-bucket algorithm. It is a no-op when rate limiting is disabled in config.
func RateLimit(cfg *config.EnvConfig) gin.HandlerFunc {
	if !cfg.RateLimitEnabled {
		return func(c *gin.Context) { c.Next() }
	}

	return RateLimitWith(cfg.RateLimitRPS, cfg.RateLimitBurst)
}

// RateLimitWith returns a per-client-IP rate-limiting middleware with an
// explicit rate (requests per second) and burst, independent of global config.
// Use it to tighten limits on sensitive routes, e.g. RateLimitPerMinute(5) for
// auth endpoints.
func RateLimitWith(rps float64, burst int) gin.HandlerFunc {
	limiter := newIPRateLimiter(rps, burst)

	return func(c *gin.Context) {
		if !limiter.limiterFor(c.ClientIP()).Allow() {
			utils.TooManyRequests(c, "")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitPerMinute returns a per-client-IP rate limiter capped at max requests
// per minute, allowing an initial burst of the same size.
func RateLimitPerMinute(max int) gin.HandlerFunc {
	return RateLimitWith(float64(max)/60.0, max)
}
