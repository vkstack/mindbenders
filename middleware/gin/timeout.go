package gin

import (
	"time"

	"github.com/gin-gonic/gin"
)

//Timeout is experimental, and not finished yet.
func Timeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)

		go func() {
			c.Next()
			finish <- struct{}{}
		}()
		<-finish
	}
}
