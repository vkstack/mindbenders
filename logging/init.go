package logging

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggeronce sync.Once
var (
	defaultLogger = &emptyLogger{}

	logger IDotpeLogger = defaultLogger
)
var host, _ = os.Hostname()

func (l *dlogger) initZap() {
	logdir := os.Getenv("LOGDIR")
	stat, err := os.Stat(logdir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal("specified path is not a directory: ", logdir)
	}
	var filename string
	var encconf zapcore.EncoderConfig = zap.NewProductionEncoderConfig()
	if os.Getenv("ENV") == "dev" {
		filename = fmt.Sprintf("app-%s-%s.log", l.app, host)
	} else {
		filename = fmt.Sprintf("app-%s.log", host)
	}
	w := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    logMazSize,
			Compress:   false,
			MaxBackups: 20,
			MaxAge:     20,
		},
	)
	encconf.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	level, err := zapcore.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		level = zap.InfoLevel
	}
	encconf.TimeKey = "time"
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encconf), w, level)
	l.zap = zap.New(core)
}

// InitLogger sets up the logger object with LoggerOptions provided.
// It returns reference logger object and error
func MustGet(opts ...Option) IDotpeLogger {
	compress, _ = strconv.ParseBool(os.Getenv("LOGCOMPRESS"))
	loggeronce.Do(func() {
		var loggr = new(dlogger)
		for _, opt := range opts {
			opt(loggr)
		}
		loggr.initZap()
		if err := loggr.finalizeEssentials(); err != nil {
			panic("unable to initialize logger")
		}
		logger = loggr
	})
	return logger
}

func DefaultLogWriter() ILogWriter { return defaultLogger }

func LogWriter() ILogWriter { return logger }
