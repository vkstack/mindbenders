package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ILogWriter interface {
	WriteLogs(context.Context, logrus.Fields, Level, string)
}

// ILogger ...
type IDotpeLogger interface {
	ILogWriter
	Gin() gin.HandlerFunc
}
