package logging

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ILogWriter interface {
	WriteLogs(context.Context, Fields, Level, string)
	Info(context.Context, Fields, string)
	Error(context.Context, Fields, string)
	Warn(context.Context, Fields, string)
	Debug(context.Context, Fields, string)
}

// ILogger ...
type IDotpeLogger interface {
	ILogWriter
	Gin() gin.HandlerFunc
}
