package logging

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"gitlab.com/dotpe/mindbenders/corel"
)

type accessLogOption func(c *gin.Context, fields Fields)

type logOption func(ctx context.Context, fields Fields)

func accessLogOptionBasic(app string) accessLogOption {
	return func(c *gin.Context, fields Fields) {
		corelid := corel.GetCorelationId(c)
		c.Writer.Header().Set("request-id", corelid.GetRequestId())
		fields["request-referer"] = c.Request.Referer()
		fields["request-clientIP"] = c.ClientIP()
		fields["request-host"] = c.Request.Host
		fields["request-method"] = c.Request.Method
		fields["request-path"] = c.FullPath()
		if fields["request-path"] != c.Request.URL.Path {
			fields["request-uripath"] = c.Request.URL.Path
		}
		fields["request-query"] = c.Request.URL.RawQuery
		fields["request-ua"] = c.Request.UserAgent()
	}
}

func AccessLogOptionRequestBody(c *gin.Context, fields Fields) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		fields["requestBody"] = string(bodyBytes)
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}

func logOptionBasic(ctx context.Context, fields Fields) {
	var corelid *corel.CoRelationId
	var err error
	if _, ok := ctx.(*gin.Context); ok {
		corelid = corel.GetCorelationId(ctx)
	} else {
		corelid, err = corel.ReadCorelId(ctx)
	}
	if err != nil {
		log.Panicln("invalid corelId: ", err.Error())
	}
	fields["sessionId"] = corelid.SessionId
	fields["requestId"] = corelid.RequestId
	fields["appRequestId"] = corelid.AppRequestId
	if len(corelid.RequestSource) != 0 {
		fields["requestSource"] = corelid.RequestSource
	}
	fields["hostname"] = host
	if os.Getenv("LOGLEVEL") == "debug" {
		fields["debug-stack-trace"] = string(debug.Stack())
	}
}
