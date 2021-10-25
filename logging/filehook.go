package logging

import (
	"fmt"
	"path"
	"time"

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
	return getFileHook(path.Join(flc.logdir, fmt.Sprintf("app-%s.log", flc.app)))
}

//absolute filename
// /home/bob/work/app.log
func getFileHook(filename string) (logrus.Hook, error) {
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}
	return rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     10,
		Level:      logrus.DebugLevel,
		Formatter:  formatter,
	})
}
