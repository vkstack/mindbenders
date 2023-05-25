package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

var (
	logMaxSize int  = 500
	compress   bool = false
)

func MustGetFileHook(app string) logrus.Hook {
	compress, _ = strconv.ParseBool(os.Getenv("LOGCOMPRESS"))
	if s, _ := strconv.Atoi(os.Getenv("LOGSIZE")); s >= 100 && s <= 1000 {
		logMaxSize = s
	}
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

func GetJSONFileHook(dir, file string) (logrus.Hook, error) {
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}

	return rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   path.Join(dir, file),
		MaxSize:    logMaxSize,
		MaxBackups: 20,
		MaxAge:     20,
		Level:      logrus.DebugLevel,
		Formatter:  formatter,
		Compress:   compress,
	})
}
