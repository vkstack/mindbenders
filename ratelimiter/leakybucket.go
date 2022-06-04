package ratelimiter

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"golang.org/x/net/context"
)

const lookupkey = "rl:locations"

/*
	It has to be used along with Postware middlewares only
*/
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

/*
	This has to be put after actual handler, since it requires input from actual handler.
	The control is passed using c.Next(), which on limit reach inform in c.c.IsAborted()
*/
func PostwareLimitUser(limiter *redis_rate.Limiter, limit redis_rate.Limit) gin.HandlerFunc {
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

/*
	This has to be put before actual handler, since it doesn't require any input from actual handler
*/
func PrewareLimitUser(limiter *redis_rate.Limiter, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.Join([]string{c.ClientIP(), c.Request.Method, c.FullPath()}, ":")
		if !Eval(c, limiter, limit, key) {
			c.Abort()
			c.JSON(http.StatusTooManyRequests, gin.H{"status": false, "message": "too many requests"})
		}
	}
}
