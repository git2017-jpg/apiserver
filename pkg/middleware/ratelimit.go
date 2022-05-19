package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/pkg/errors"
	"monitor-apiserver/pkg/response"
)

// ErrLimitExceed is service unavailable due to rate limit exceeded.
var ErrLimitExceed = errors.New("[RATELIMIT] service unavailable due to rate limit exceeded")

// LimiterOption is ratelimit option.
type LimiterOption func(*limiterOptions)

// WithLimiter set Limiter implementation,
// default is bbr limiter
func WithLimiter(limiter ratelimit.Limiter) LimiterOption {
	return func(o *limiterOptions) {
		o.limiter = limiter
	}
}

type limiterOptions struct {
	limiter ratelimit.Limiter
}

// Ratelimit a circuit breaker middleware
func Ratelimit(opts ...LimiterOption) gin.HandlerFunc {
	options := &limiterOptions{
		limiter: bbr.NewLimiter(),
	}
	for _, o := range opts {
		o(options)
	}
	return func(c *gin.Context) {
		done, e := options.limiter.Allow()
		if e != nil {
			// rejected
			response.JSON(c, ErrLimitExceed, nil)
			c.Abort()
			return
		}
		// allowed
		done(ratelimit.DoneInfo{Err: c.Request.Context().Err()})

		c.Next()
	}
}
