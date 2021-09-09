package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type FileLogConfig struct {
	logdir, app string
}

func NewFileLogConfig(logdir, app string) ILogConfig {
	return &FileLogConfig{
		logdir: logdir,
		app:    app,
	}
}

func (flc *FileLogConfig) getHook() (logrus.Hook, error) {
	formatter := &logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}
	return rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logfile.log",
		MaxSize:    5,
		MaxBackups: 7,
		MaxAge:     7,
		Level:      logrus.DebugLevel,
		Formatter:  formatter,
	})
}
