package logging

import (
	"fmt"
	"log"
	"os"

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

func (l *dlogger) initZap() {
	logdir := os.Getenv("LOGDIR")
	stat, err := os.Stat(logdir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("specified path is not a directory: ", logdir)
	}
	var filename string
	var zencconf zapcore.EncoderConfig
	if os.Getenv("ENV") == "dev" {
		zencconf = zap.NewDevelopmentEncoderConfig()
		filename = fmt.Sprintf("app-%s-%s.log", l.app, host)
	} else {
		zencconf = zap.NewProductionEncoderConfig()
		filename = fmt.Sprintf("app-%s.log", host)
	}
	level, err := zapcore.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		level = zap.InfoLevel
	}
	zencconf.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zencconf.TimeKey = "time"
	var zsyncer = zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
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
	l.zap = zap.New(core)
}
