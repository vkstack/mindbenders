package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	PanicLevel  Level = Level(zap.PanicLevel)
	FatalLevel  Level = Level(zap.FatalLevel)
	ErrorLevel  Level = Level(zap.ErrorLevel)
	WarnLevel   Level = Level(zap.WarnLevel)
	InfoLevel   Level = Level(zap.InfoLevel)
	DebugLevel  Level = Level(zap.DebugLevel)
	DPanicLevel Level = Level(zap.DPanicLevel)
)

func (l Level) String() string { return zapcore.Level(l).String() }
