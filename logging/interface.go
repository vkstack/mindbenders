package logging

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ILogWriter interface {
	WriteLogs(context.Context, Fields, Level, string)
}

// ILogger ...
type IDotpeLogger interface {
	ILogWriter
	Gin() gin.HandlerFunc
}
