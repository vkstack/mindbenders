package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type dlogger struct {
	app, appId, env,
	wd string // Working directory of the application

	logger   *logrus.Logger
	accopts  []accessLogOption
	loptions []logOption
}

func (dlogger *dlogger) safeRunLogOptions(ctx context.Context, fields *logrus.Fields) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("unknown error while operating logOptions", r)
		}
	}()
	for _, opt := range dlogger.loptions {
		opt(ctx, fields)
	}
}

func (dlogger *dlogger) safeRunAccessLogOptions(c *gin.Context, fields *logrus.Fields) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("unknown error while operating accesslogOptions", r)
		}
	}()
	for _, opt := range dlogger.accopts {
		opt(c, fields)
	}
}

func (dlogger *dlogger) finalizeEssentials() error {
	if dlogger.logger == nil || dlogger.logger.Hooks == nil {
		hook, err := getFileHook("app.log")
		if err != nil {
			return err
		}
		WithHook(hook)(dlogger)
	}
	if dlogger.loptions == nil {
		dlogger.loptions = append(dlogger.loptions, logOptionBasic)
	}
	if dlogger.accopts == nil {
		dlogger.accopts = append(dlogger.accopts, accessLogOptionBasic(dlogger.app))
	}
	return nil
}

//WriteLogs writes log
func (dLogger *dlogger) WriteLogs(ctx context.Context, fields logrus.Fields, cb logrus.Level, MessageKey string) {
	if ctx == nil {
		return
	}
	if len(dLogger.appId) > 0 {
		fields["appID"] = dLogger.appId
	}
	dLogger.safeRunLogOptions(ctx, &fields)
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
	fields["caller"] = strings.ReplaceAll(fields["caller"].(string), dLogger.wd, "")
	entry := dLogger.logger.WithFields(fields)
	if t, ok := ctx.Value("time").(time.Time); ok {
		entry.Time = t
	} else {
		entry.Time = time.Now()
	}
	entry.Log(cb, MessageKey)
}

//GinLogger returns a gin.HandlerFunc middleware
func (dLogger *dlogger) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("time", start)
		var fields = logrus.Fields{}
		dLogger.safeRunAccessLogOptions(c, &fields)
		var level = new(logrus.Level)
		*level = logrus.InfoLevel

		//deferred request log
		defer dLogger.WriteLogs(c, fields, *level, "access-log")

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
