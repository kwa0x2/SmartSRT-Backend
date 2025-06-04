package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP() + ":" + ctx.Request.URL.Path

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		var validRequests []time.Time
		for _, t := range rl.requests[key] {
			if now.Sub(t) < rl.window {
				validRequests = append(validRequests, t)
			}
		}
		rl.requests[key] = validRequests

		if len(rl.requests[key]) >= rl.limit {
			ctx.JSON(http.StatusTooManyRequests, utils.NewMessageResponse("You're sending too many requests. Please slow down."))
			ctx.Abort()
			return
		}

		rl.requests[key] = append(rl.requests[key], now)

		ctx.Next()
	}
}
