package logging

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var loggeronce sync.Once
var (
	defaultLogger = &emptyLogger{}

	logger IDotpeLogger = defaultLogger
)
var host, _ = os.Hostname()

// InitLogger sets up the logger object with LoggerOptions provided.
// It returns reference logger object and error
func MustGet(opts ...Option) IDotpeLogger {
	loggeronce.Do(func() {
		var loggr = new(dlogger)
		loggr.initExporter()
		for _, opt := range opts {
			opt(loggr)
		}
		if err := loggr.finalizeEssentials(); err != nil {
			panic("unable to initialize logger")
		}
		logger = loggr
		level, err := logrus.ParseLevel(os.Getenv("LOGLEVEL"))
		if err != nil {
			level = logrus.InfoLevel
		}
		loggr.logger.SetLevel(level)
	})
	return logger
}

func DefaultLogWriter() ILogWriter { return defaultLogger }

func LogWriter() ILogWriter { return logger }
