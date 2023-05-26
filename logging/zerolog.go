package logging

import (
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getZero(app string) *zerolog.Logger {
	rotator := zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   path.Join(getlogdir(), getLogFileName(app)),
		MaxSize:    logMaxSize,
		Compress:   false,
		MaxBackups: 20,
		MaxAge:     20,
	})
	var multi zerolog.LevelWriter
	if os.Getenv("ENV") == "dev" {
		multi = zerolog.MultiLevelWriter(rotator, os.Stdout)
	} else {
		multi = zerolog.MultiLevelWriter(rotator)
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(multi)
	return &logger
}

func (dLogger *dlogger) zeroWrite(fields Fields, cb Level, MessageKey string) {
	zerolevel := mapZerolog[cb]
	event := dLogger.zero.WithLevel(zerolevel)
	event.Time("time", time.Now())
	if t, ok := fields["time"]; ok {
		if ts, ok := t.(time.Time); ok {
			event.Time("time", ts)
		}
		delete(fields, "time")
	}
	for k, v := range fields {
		event = event.Any(k, v)
	}
	event.Msg(MessageKey)
}
