package logging

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

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
	var fsize int64
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the io.ReadCloser to its original state
		if strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
			fsize = fileSize(*c.Request) // find file size
		}
		if fsize == 0 {
			fields["request-body"] = string(bodyBytes)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the io.ReadCloser to its original state
		fields["request-fileLength"] = fsize
	}
}

func fileSize(req http.Request) int64 {
	if err := req.ParseMultipartForm(maxMultiPartSize); err != nil {
		log.Panicln("multipart parse issue : ", err.Error())
	}
	var fsize int64
	for _, files := range req.MultipartForm.File {
		for _, file := range files {
			fsize += file.Size
		}
	}
	return fsize
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
