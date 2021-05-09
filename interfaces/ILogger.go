package interfaces

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//ILogger ...
type IDotpeLogger interface {
	WriteLogs(context.Context, logrus.Fields, logrus.Level, string, ...interface{})
	GinLogger() gin.HandlerFunc
}
