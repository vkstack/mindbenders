package logging

import (
	"github.com/rs/zerolog"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	PanicLevel Level = Level(zerolog.PanicLevel)
	FatalLevel Level = Level(zerolog.FatalLevel)
	ErrorLevel Level = Level(zerolog.ErrorLevel)
	WarnLevel  Level = Level(zerolog.WarnLevel)
	InfoLevel  Level = Level(zerolog.InfoLevel)
	DebugLevel Level = Level(zerolog.DebugLevel)
	// DPanicLevel Level = Level(zerolog.DPanicLevel)
)

func (l Level) String() string { return zerolog.Level(l).String() }
