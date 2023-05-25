package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type emptyLogger struct{}

func (el *emptyLogger) WriteLogs(context.Context, logrus.Fields, Level, string) {}

func (el *emptyLogger) Gin() gin.HandlerFunc { return func(c *gin.Context) {} }
