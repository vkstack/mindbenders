package logging

import (
	"context"

	"github.com/gin-gonic/gin"
)

type emptyLogger struct{}

// Deprecated: Use [Log] instead.
func (el *emptyLogger) WriteLogs(context.Context, Fields, Level, string) {}

func (el *emptyLogger) Log(context.Context, Fields, Level, string) {}

func (el *emptyLogger) Info(context.Context, Fields, string)  {}
func (el *emptyLogger) Error(context.Context, Fields, string) {}
func (el *emptyLogger) Warn(context.Context, Fields, string)  {}
func (el *emptyLogger) Debug(context.Context, Fields, string) {}
func (el *emptyLogger) Fatal(context.Context, Fields, string) {}
func (el *emptyLogger) Panic(context.Context, Fields, string) {}

func (el *emptyLogger) Gin() gin.HandlerFunc { return func(c *gin.Context) {} }
