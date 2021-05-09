package gin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/interfaces"
)

//Recovery returns a gin.HandlerFunc having recovery solution
func Recovery(l interfaces.IDotpeLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := fmt.Sprintf("%v\n%s", err, debug.Stack())
				//copied from /usr/local/go/src/runtime/debug/stack.go | gin@v1.6.3
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							c.Error(err.(error)) // nolint: errcheck
							c.Abort()
							l.WriteLogs(c, logrus.Fields{"stacktrace": stack}, logrus.TraceLevel, "BrokenPipe")
							return
						}
					}
				}
				l.WriteLogs(c, logrus.Fields{"stacktrace": stack}, logrus.FatalLevel, "Panic")
				c.JSON(http.StatusExpectationFailed, map[string]interface{}{
					"message": "something went wrong",
					"status":  false,
				})
			}
		}()
		c.Next()
	}
}
