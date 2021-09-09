package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/dotpe/mindbenders/corel"
)

type dlogger struct {
	Lops   LoggerOptions
	Logger *logrus.Logger
}

//KibanaConfig Mandatory for kibana logging

type ILogConfig interface {
	getHook() (logrus.Hook, error)
}

// type

//LoggerOptions is set of config data for logg
type LoggerOptions struct {
	iconfig ILogConfig
	APP,    // Service
	APPID, // Service application ID
	LOGENV, // Dev/Debug/Production
	Hostname,
	WD string // Working directory of the application
	DisableJSONLogging bool
}

//WriteLogs writes log
func (dLogger *dlogger) WriteLogs(ctx context.Context, fields logrus.Fields, cb logrus.Level, MessageKey string) {
	if ctx == nil {
		return
	}

	for idx := range fields {
		switch fields[idx].(type) {
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64,
			string, bool:
		default:
			tmp, _ := json.Marshal(fields[idx])
			fields[idx] = string(tmp)
		}
	}
	if _, ok := fields["caller"]; !ok {
		pc, file, line, _ := runtime.Caller(1)
		_, funcname := filepath.Split(runtime.FuncForPC(pc).Name())
		file = strings.Trim(file, " ")
		funcname = strings.Trim(funcname, " ")
		fields["caller"] = fmt.Sprintf("%s:%d\n%s", file, line, funcname)
	}
	fields["caller"] = strings.ReplaceAll(fields["caller"].(string), dLogger.Lops.WD, "")
	fields["appID"] = dLogger.Lops.APPID
	coRelationID, _ := corel.GetCorelationId(ctx)
	fields["requestID"] = coRelationID.RequestID
	fields["sessionID"] = coRelationID.SessionID
	fields["hop"] = coRelationID.Hop
	if coRelationID.OriginApp != "" {
		fields["OriginApp"] = coRelationID.OriginApp
		fields["OriginHost"] = coRelationID.OriginHost
	}
	entry := dLogger.Logger.WithFields(fields)
	entry.Log(cb, MessageKey)
}

//GinLogger returns a gin.HandlerFunc middleware
func (dLogger *dlogger) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		start := time.Now()
		var corelid corel.CoRelationId
		c.ShouldBindHeader(&corelid)
		corelid.OnceMust(c, dLogger.Lops.APP)
		corel.GinSetCoRelID(c, &corelid)
		fields := logrus.Fields{
			"referer":     c.Request.Referer(),
			"clientIP":    c.ClientIP(),
			"host":        c.Request.Host,
			"hostname":    dLogger.Lops.Hostname,
			"method":      c.Request.Method,
			"path":        c.FullPath(),
			"uriparams":   parseGinUriParams(c.Params),
			"queryparams": c.Request.URL.Query(),
			"requestID":   corelid.RequestID,
			"sessionID":   corelid.SessionID,
			"hop":         corelid.Hop,
			"userAgent":   c.Request.UserAgent(),
		}

		var level = new(logrus.Level)
		*level = logrus.InfoLevel

		//deferred request log
		defer dLogger.WriteLogs(c, fields, *level, "access-log")
		var bodyBytes []byte
		if c.Request.Body != nil && !dLogger.Lops.DisableJSONLogging {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			fields["requestBody"] = string(bodyBytes)
			// Restore the io.ReadCloser to its original state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		fields["statusCode"] = 0
		c.Next()
		stop := time.Since(start)
		fields["latency"] = int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		code := c.Writer.Status()

		fields["statusCode"] = code
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		fields["dataLength"] = dataLength

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			*level = logrus.ErrorLevel
		} else if code > 499 {
			*level = logrus.ErrorLevel
		} else if code > 399 {
			*level = logrus.WarnLevel
		}
	}
}

func parseGinUriParams(params gin.Params) map[string]interface{} {
	parsedParams := make(map[string]interface{})
	for _, p := range params {
		parsedParams[p.Key] = p.Value
	}
	return parsedParams
}
