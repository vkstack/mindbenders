package logging

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/corel"
)

type accessLogOption func(c *gin.Context, fields *logrus.Fields)

func accessLogOptionBasic(app string) accessLogOption {
	return func(c *gin.Context, fields *logrus.Fields) {
		corelid, _ := corel.GetCorelationId(c)
		c.Writer.Header().Set("request-id", corelid.RequestID)
		(*fields)["referer"] = c.Request.Referer()
		(*fields)["clientIP"] = c.ClientIP()
		(*fields)["host"] = c.Request.Host
		(*fields)["method"] = c.Request.Method
		(*fields)["path"] = c.FullPath()
		(*fields)["uriparams"] = parseGinUriParams(c.Params)
		(*fields)["queryparams"] = c.Request.URL.Query()
		(*fields)["userAgent"] = c.Request.UserAgent()
	}
}

func parseGinUriParams(params gin.Params) map[string]interface{} {
	parsedParams := make(map[string]interface{})
	for _, p := range params {
		parsedParams[p.Key] = p.Value
	}
	return parsedParams
}

func AccessLogOptionRequestBody(c *gin.Context, fields *logrus.Fields) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		(*fields)["requestBody"] = string(bodyBytes)
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}

func AccessLogOptionJWTLogging(c *gin.Context, fields *logrus.Fields) {
	corelid, _ := corel.GetCorelationId(c)
	if len(corelid.Auth) > 0 {
		(*fields)["token"] = corelid.Auth
	}
}

type logOption func(ctx context.Context, fields *logrus.Fields)

func logOptionBasic(ctx context.Context, fields *logrus.Fields) {
	coRelationID, err := corel.GetCorelationId(ctx)
	if err != nil {
		log.Panicln("invalid corelId: ", err.Error())
	}
	(*fields)["requestID"] = coRelationID.RequestID
	(*fields)["sessionID"] = coRelationID.SessionID
	(*fields)["hostname"] = host
	if coRelationID.OriginApp != "" {
		(*fields)["OriginApp"] = coRelationID.OriginApp
		(*fields)["OriginHost"] = coRelationID.OriginHost
	}
}
