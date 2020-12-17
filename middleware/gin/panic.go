package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recovery(l logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				_ctx, _ := c.Get("context")
				ctx := _ctx.(context.Context)
				l.WriteLogs(ctx, logrus.Fields{
					"params": fmt.Sprintf("%v\n%s", r, debug.Stack()),
				}, logrus.FatalLevel, "Panic")
				c.JSON(http.StatusExpectationFailed, map[string]interface{}{
					"message": "something went wrong",
					"status":  false,
				})
			}
		}()
		c.Next()
	}
}
