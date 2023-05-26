package logging

import (
	"fmt"

	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
	// DPanicLevel Level = Level(zerolog.DPanicLevel)
)

func (l Level) String() string {
	if str, ok := mapLevelStr[l]; ok {
		return str
	}
	return fmt.Sprintf("invalid:%d", l)
}

var (
	mapZerolog = map[Level]zerolog.Level{
		PanicLevel: zerolog.PanicLevel,
		FatalLevel: zerolog.FatalLevel,
		ErrorLevel: zerolog.ErrorLevel,
		WarnLevel:  zerolog.WarnLevel,
		InfoLevel:  zerolog.InfoLevel,
		DebugLevel: zerolog.DebugLevel,
	}

	mapZap = map[Level]zapcore.Level{
		PanicLevel: zapcore.PanicLevel,
		FatalLevel: zapcore.FatalLevel,
		ErrorLevel: zapcore.ErrorLevel,
		WarnLevel:  zapcore.WarnLevel,
		InfoLevel:  zapcore.InfoLevel,
		DebugLevel: zapcore.DebugLevel,
	}

	mapLevelStr = map[Level]string{
		PanicLevel: "panic",
		FatalLevel: "fatal",
		ErrorLevel: "error",
		WarnLevel:  "warn",
		InfoLevel:  "info",
		DebugLevel: "debug",
		TraceLevel: "trace",
	}
)
