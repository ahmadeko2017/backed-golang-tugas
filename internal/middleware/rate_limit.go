package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

var limiter = NewIPRateLimiter(5, 10) // 5 req/s, burst 10

// RateLimitMiddleware limits requests per IP
func RateLimitMiddleware(rps int) gin.HandlerFunc {
	// Override default if needed, though here using global variable for simplicity in this task
	// Ideally we'd pass rps to NewIPRateLimiter
	// For now keeping interface same: RateLimitMiddleware(rps int)

	// We update the limiter with new RPS
	limiter = NewIPRateLimiter(rate.Limit(rps), rps*2)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
