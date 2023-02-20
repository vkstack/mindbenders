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

func MustGetFileHook(app string) logrus.Hook {
	logdir := os.Getenv("LOGDIR")
	stat, err := os.Stat(logdir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("specified path is not a directory: ", logdir)
	}
	hn, _ := os.Hostname()
	var filename string
	if os.Getenv("ENV") == "dev" {
		filename = fmt.Sprintf("app-%s-%s.log", app, hn)
	} else {
		filename = fmt.Sprintf("app-%s.log", hn)
	}
	hook, err := GetJSONFileHook(logdir, filename)
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
		MaxSize:    1024,
		MaxBackups: 20,
		MaxAge:     20,
		Level:      logrus.DebugLevel,
		Formatter:  formatter,
		Compress:   true,
	})
}
