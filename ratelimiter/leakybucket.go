package ratelimiter

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"golang.org/x/net/context"
)

const lookupkey = "rl:locations"

func IsRunnable(c *gin.Context, resource string) bool {
	c.Set(lookupkey, resource)
	c.Next()
	return !c.IsAborted()
}

func Eval(ctx context.Context, limiter *redis_rate.Limiter, limit redis_rate.Limit, key string) bool {
	res, err := limiter.Allow(ctx, key, limit)
	if err != nil {
		return false
	}
	return res.Allowed+res.Remaining > 0
}

func LimitUserByResourceExt(limiter *redis_rate.Limiter, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		if data, ok := c.Get(lookupkey); ok {
			if resource, ok := data.(string); ok {
				if !Eval(c, limiter, limit, resource) {
					c.Abort()
					c.JSON(http.StatusTooManyRequests, gin.H{"status": false, "message": "too many requests"})
				}
			}
		}
	}
}

func LimitUserByResource(limiter *redis_rate.Limiter, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.Join([]string{c.ClientIP(), c.Request.Method, c.FullPath()}, ":")
		if !Eval(c, limiter, limit, key) {
			c.Abort()
			c.JSON(http.StatusTooManyRequests, gin.H{"status": false, "message": "too many requests"})
		}
	}
}
