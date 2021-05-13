package interfaces

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ILogWriter interface {
	WriteLogs(context.Context, logrus.Fields, logrus.Level, string, ...interface{})
}

//ILogger ...
type IDotpeLogger interface {
	ILogWriter
	GinLogger() gin.HandlerFunc
}
