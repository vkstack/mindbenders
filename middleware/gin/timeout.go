package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)

		go func() {
			c.Next()
			finish <- struct{}{}
		}()

		select {
		// case <-time.After(t):
		// 	c.JSON(504, "timeout")
		// 	c.Abort()
		case <-finish:
		}
	}
}