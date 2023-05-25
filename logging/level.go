package logging

import (
	"github.com/sirupsen/logrus"
)

type Level int8

type Fields map[string]interface{}

const (
	PanicLevel Level = Level(logrus.PanicLevel)
	FatalLevel Level = Level(logrus.FatalLevel)
	ErrorLevel Level = Level(logrus.ErrorLevel)
	WarnLevel  Level = Level(logrus.WarnLevel)
	InfoLevel  Level = Level(logrus.InfoLevel)
	DebugLevel Level = Level(logrus.DebugLevel)
)

func (l Level) String() string { return logrus.Level(l).String() }
