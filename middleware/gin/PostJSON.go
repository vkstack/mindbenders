package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/interfaces"
)

// var logger func(context.Context, logrus.Fields, logrus.Level, string)

func PostJSONValidator(l interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var byteData []byte
		if c.Request.Body != nil {
			byteData, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteData))
		var jsonData interface{}
		err := json.Unmarshal(byteData, &jsonData)
		ictx, _ := c.Get("context")
		ctx := ictx.(context.Context)
		fields := logrus.Fields{
			"input":    string(byteData),
			"clientIP": c.ClientIP(),
			"path":     c.Request.URL.Path,
		}
		if err != nil {
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
