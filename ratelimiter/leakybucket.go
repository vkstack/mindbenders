package ratelimiter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
)

func LimitUserByResource(limiter *redis_rate.Limiter, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.Join([]string{c.ClientIP(), c.Request.Method, c.FullPath()}, ":")
		res, err := limiter.Allow(c, key, limit)
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.Allowed+res.Remaining <= 0 {
			c.Abort()
			c.JSON(http.StatusTooManyRequests, gin.H{"status": false, "message": "too many requests"})
		}
	}
}

// var once sync.Once
// var limiter *redis_rate.Limiter

// type _rediser interface {
// 	*redis.Client | *redis.ClusterClient
// 	Del(ctx context.Context, keys ...string) *redis.IntCmd
// }

// func NewDistLimiter[R _rediser](rds R) *redis_rate.Limiter {
// 	// r, _ := rds.(*redis.Client)
// 	redis_rate.NewLimiter(rds)
// }

// func NewDistLimiterWithRedisCluster(rds *redis.ClusterClient) *redis_rate.Limiter {
// 	once.Do(
// 		func() {
// 			limiter = redis_rate.NewLimiter(rds)
// 		},
// 	)
// 	return limiter
// }

// func NewDistLimiterWithRedis(rds *redis.Client) *redis_rate.Limiter {
// 	once.Do(
// 		func() {
// 			limiter = redis_rate.NewLimiter(rds)
// 		},
// 	)
// 	return limiter
// }

// func Limiter(ctx context.Context, key string, lim redis_rate.Limit) bool {
// 	res, err := limiter.Allow(ctx, key, lim)
// 	// res, err := limiter.AllowN(ctx, key, lim, lim.Rate)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	fmt.Print("limit: ", res.Limit, " allowed: ", res.Allowed, " remaining: ", res.Remaining, " time: ", time.Now().Unix(), "\t")
// 	return res.Allowed+res.Remaining > 0
// }
