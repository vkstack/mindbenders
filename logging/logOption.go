package logging

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/corel"
)

const (
	maxMultiPartSize = 100 << 20
)

type accessLogOption func(c *gin.Context, fields logrus.Fields)

func accessLogOptionBasic(app string) accessLogOption {
	return func(c *gin.Context, fields logrus.Fields) {
		corelid, _ := corel.GetCorelationId(c)
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

func AccessLogOptionRequestBody(c *gin.Context, fields logrus.Fields) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		dcreq := c.Request.Clone(c)
		fmt.Println("request body: ", dcreq)
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		fields["request-body"] = string(bodyBytes)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the io.ReadCloser to its original state
	}
}

type logOption func(ctx context.Context, fields logrus.Fields)

func logOptionBasic(ctx context.Context, fields logrus.Fields) {
	coRelationID, err := corel.GetCorelationId(ctx)
	if err != nil {
		log.Panicln("invalid corelId: ", err.Error())
	}
	coRelationID.Logrus(fields)
	fields["hostname"] = host
	if os.Getenv("LOGLEVEL") == "debug" {
		fields["debug-stack-trace"] = string(debug.Stack())
	}
}
