package uiware

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var mandatorilyAllowedHeaders = hashset.New("Session_id", "Authorization")
var additionalHeaders []string

var DefaultAllowedHeaders = []string{
	"app_version", "device", "device_id",
}

func AllowCorsHeaders(r *gin.Engine, headers ...string) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	for _, val := range mandatorilyAllowedHeaders.Values() {
		headers = append(headers, val.(string))
	}
	additionalHeaders = headers
	config.AddAllowHeaders(headers...)
	r.Use(cors.New(config))
}

func HeaderLogoption(c *gin.Context, filelds *logrus.Fields) {
	if filelds == nil || c == nil {
		return
	}
	for _, k := range additionalHeaders {
		if hv := c.GetHeader(k); hv != "" {
			(*filelds)[k] = hv
		}
	}
}
