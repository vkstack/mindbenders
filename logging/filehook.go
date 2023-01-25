package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func MustGetHook(app string) logrus.Hook {
	logdir := os.Getenv("LOGDIR")
	stat, err := os.Stat(logdir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("specified path is not a directory: ", logdir)
	}
	hook, err := GetJSONFileHook(logdir, fmt.Sprintf("app-%s.log", app))
	if err != nil {
		log.Fatalf("unable to get file hook:%v\n", err)
	}
	return hook
}

// absolute filename
// /home/bob/work/app.log
func GetJSONFileHook(dir, file string) (logrus.Hook, error) {
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}
	return rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   path.Join(dir, file),
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     10,
		Level:      logrus.DebugLevel,
		Formatter:  formatter,
	})
}
