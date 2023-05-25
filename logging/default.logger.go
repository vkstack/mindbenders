package logging

import (
	"context"

	"github.com/gin-gonic/gin"
)

type emptyLogger struct{}

func (el *emptyLogger) WriteLogs(context.Context, Fields, Level, string) {}

func (el *emptyLogger) Gin() gin.HandlerFunc { return func(c *gin.Context) {} }
