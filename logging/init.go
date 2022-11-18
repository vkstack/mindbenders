package logging

import (
	"os"
	"sync"
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
		for _, opt := range opts {
			opt(loggr)
		}
		if err := loggr.finalizeEssentials(); err != nil {
			panic("unable to initialize logger")
		}
		logger = loggr
	})
	return logger
}

func DefaultLogWriter() ILogWriter { return defaultLogger }

func LogWriter() ILogWriter { return logger }
