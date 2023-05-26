package logging

import (
	"math"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger returns a gin.HandlerFunc middleware
func (dLogger *dlogger) Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var fields = Fields{}
		dLogger.safeRunAccessLogOptions(c, fields)
		var level Level
		level = InfoLevel

		//deferred request log
		fields["time"] = start
		defer dLogger.WriteLogs(c, fields, level, "access-log")

		fields["request-statusCode"] = 0
		c.Next()
		stop := time.Since(start)
		fields["request-latency"] = int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		code := c.Writer.Status()

		fields["request-statusCode"] = code
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		fields["request-dataLength"] = dataLength

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			level = ErrorLevel
		} else if code > 499 {
			level = ErrorLevel
		} else if code > 399 {
			level = WarnLevel
		}
	}
}
