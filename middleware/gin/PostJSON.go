package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// var logger func(context.Context, logrus.Fields, logrus.Level, string)

type logger interface {
	WriteLogs(context.Context, logrus.Fields, logrus.Level, string, ...interface{})
}

func PostJSONValidator(l logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonData interface{}
		byteData, _ := c.GetRawData()
		err := json.Unmarshal(byteData, &jsonData)
		ictx, _ := c.Get("context")
		ctx := ictx.(context.Context)
		fields := logrus.Fields{
			"input":    string(byteData),
			"clientIP": c.ClientIP(),
			"path":     c.Request.URL.Path,
		}
		if err != nil {
			fields["errorMsg"] = error.Error(err)
			l.WriteLogs(ctx, fields, logrus.ErrorLevel, "bad_json")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"error":   error.Error(err),
				"message": "Bad JSON.",
			})
			c.Abort()
			return
		}
		l.WriteLogs(ctx, fields, logrus.InfoLevel, "request_json")
		c.Set("jsonByte", byteData)
	}
}
