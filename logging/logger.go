package logging

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

const logMaxSize int = 500

var (
	compress bool = false
)

type Fields map[string]interface{}

type dlogger struct {
	app,
	appId,
	env string

	zap   *zap.Logger
	iszap bool

	zero     *zerolog.Logger
	writer   func(Fields, Level, string)
	accopts  []accessLogOption
	loptions []logOption

	fieldexecutor         []func(Fields)
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
	if dlogger.loptions == nil {
		dlogger.loptions = append(dlogger.loptions, logOptionBasic)
	}
	if dlogger.accopts == nil {
		dlogger.accopts = append(dlogger.accopts, accessLogOptionBasic(dlogger.app))
	}
	return nil
}

func (dLogger *dlogger) write(ctx context.Context, fields Fields, cb Level, MessageKey string) {
	if ctx == nil {
		return
	}
	for _, ex := range dLogger.fieldexecutor {
		ex(fields)
	}
	dLogger.addMetrics(cb, fields["caller"].(string))
	dLogger.safeRunLogOptions(ctx, fields)
	dLogger.writer(fields, cb, MessageKey)
}

// WriteLogs writes log
func (dLogger *dlogger) WriteLogs(ctx context.Context, fields Fields, cb Level, MessageKey string) {
	dLogger.write(ctx, fields, cb, MessageKey)
}

func (dLogger *dlogger) Info(ctx context.Context, fields Fields, MessageKey string) {
	dLogger.write(ctx, fields, InfoLevel, MessageKey)
}

func (dLogger *dlogger) Error(ctx context.Context, fields Fields, MessageKey string) {
	dLogger.write(ctx, fields, ErrorLevel, MessageKey)
}

func (dLogger *dlogger) Warn(ctx context.Context, fields Fields, MessageKey string) {
	dLogger.write(ctx, fields, WarnLevel, MessageKey)
}

func (dLogger *dlogger) Debug(ctx context.Context, fields Fields, MessageKey string) {
	dLogger.write(ctx, fields, DebugLevel, MessageKey)
}
