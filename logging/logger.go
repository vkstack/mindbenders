package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type dlogger struct {
	app,
	appId,
	env string

	logger   *logrus.Logger
	accopts  []accessLogOption
	loptions []logOption

	metricCollectionLevel Level
	collector             *prometheus.CounterVec
}

func (dlogger *dlogger) safeRunLogOptions(ctx context.Context, fields Fields) {
	for _, opt := range dlogger.loptions {
		if opt != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						stack := fmt.Sprintf("%v\n%s", r, debug.Stack())
						log.Println("unknown error while operating logOptions\n", stack)
					}
				}()
				opt(ctx, fields)
			}()
		}
	}
}

func (dlogger *dlogger) safeRunAccessLogOptions(c *gin.Context, fields Fields) {
	defer func() {
		if r := recover(); r != nil {
			stack := fmt.Sprintf("%v\n%s", r, debug.Stack())
			log.Println("unknown error while operating accesslogOptions\n", stack)
		}
	}()
	for _, opt := range dlogger.accopts {
		if opt != nil {
			opt(c, fields)
		}
	}
}

func (dlogger *dlogger) finalizeEssentials() error {
	if dlogger.logger == nil || dlogger.logger.Hooks == nil {
		hook, err := GetJSONFileHook(".", "app.log")
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

// WriteLogs writes log
func (dLogger *dlogger) WriteLogs(ctx context.Context, fields Fields, cb Level, MessageKey string) {
	if ctx == nil {
		return
	}
	if len(dLogger.appId) > 0 {
		fields["appID"] = dLogger.appId
	}
	dLogger.safeRunLogOptions(ctx, fields)
	for idx := range fields {
		switch x := fields[idx].(type) {
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64,
			string, bool:
		case fmt.Stringer:
			fields[idx] = x.String()
		case error:
			fields[idx] = x.Error()
		default:
			tmp, _ := json.Marshal(fields[idx])
			fields[idx] = string(tmp)
		}
	}
	pc, file, line, _ := runtime.Caller(1)
	_, funcname := filepath.Split(runtime.FuncForPC(pc).Name())
	file = canonicalFile(strings.Trim(file, "/"))
	funcname = strings.Trim(funcname, " ")
	fields["caller"] = fmt.Sprintf("%s:%d\n%s", file, line, funcname)
	dLogger.addMetrics(cb, MessageKey, fmt.Sprintf("%s:%d", file, line))
	entry := dLogger.logger.WithFields(logrus.Fields(fields))
	entry.Time = time.Now()
	if t, ok := fields["time"]; ok {
		if ts, ok := t.(time.Time); ok {
			entry.Time = ts
		}
		delete(fields, "time")
	}
	entry.Log(logrus.Level(cb), MessageKey)
}

func canonicalFile(file string) string {
	file = strings.Trim(file, "/")
	parts := strings.Split(file, "/")
	return strings.Join(parts[:len(parts)/3], "/") +
		"\n" +
		strings.Join(parts[len(parts)/3:], "/")
}

// GinLogger returns a gin.HandlerFunc middleware
func (dLogger *dlogger) Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var fields = Fields{}
		dLogger.safeRunAccessLogOptions(c, fields)
		var level Level
		level = InfoLevel

		//deferred request log
		fields["time"] = start
		defer dLogger.WriteLogs(c, fields, level, "access-log")

		fields["request-statusCode"] = 0
		c.Next()
		stop := time.Since(start)
		fields["request-latency"] = int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		code := c.Writer.Status()

		fields["request-statusCode"] = code
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		fields["request-dataLength"] = dataLength

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			level = ErrorLevel
		} else if code > 499 {
			level = ErrorLevel
		} else if code > 399 {
			level = WarnLevel
		}
	}
}
