package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func (l *dlogger) enzap(fields Fields) (zfields []zap.Field) {
	for k, v := range fields {
		zfields = append(zfields, zap.Any(k, v))
	}
	return
}
func getlogdir() string {
	logdir := os.Getenv("LOGDIR")
	stat, err := os.Stat(logdir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("specified path is not a directory: ", logdir)
	}
	return logdir
}

func getLogFileName(app string) string {
	if os.Getenv("ENV") == "dev" {
		return fmt.Sprintf("app-%s.log", app)
	}
	return fmt.Sprintf("app-%s.log", host)
}

func getZap(app string) *zap.Logger {
	var zencconf zapcore.EncoderConfig
	if os.Getenv("ENV") == "dev" {
		zencconf = zap.NewDevelopmentEncoderConfig()
	} else {
		zencconf = zap.NewProductionEncoderConfig()
	}
	level, err := zapcore.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		level = zap.InfoLevel
	}
	zencconf.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zencconf.TimeKey = "time"
	var zsyncer = zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(getlogdir(), getLogFileName(app)),
		MaxSize:    logMaxSize,
		Compress:   false,
		MaxBackups: 20,
		MaxAge:     20,
	})
	if os.Getenv("ENV") == "dev" {
		zsyncer = zapcore.NewMultiWriteSyncer(zsyncer, os.Stdout)
	}
	var core zapcore.Core = zapcore.NewCore(
		zapcore.NewJSONEncoder(zencconf),
		zsyncer,
		level,
	)
	return zap.New(core)
}

func (dLogger *dlogger) zapWrite(fields Fields, cb Level, MessageKey string) {
	zlevel := mapZap[cb]
	entry := dLogger.zap.Check(zlevel, MessageKey)
	if t, ok := fields["time"]; ok {
		if ts, ok := t.(time.Time); ok {
			entry.Time = ts
		}
		delete(fields, "time")
	}
	zfields := dLogger.enzap(fields)
	entry.Write(zfields...)
}
