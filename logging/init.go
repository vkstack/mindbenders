package logging

import (
	"os"
	"strconv"
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
func getlogger(opts ...Option) *dlogger {
	var loggr = new(dlogger)
	for _, opt := range opts {
		opt(loggr)
	}
	loggr.fieldexecutor = append(loggr.fieldexecutor, fieldNormalize, caller, checkReservedKeys)
	if loggr.iszap {
		loggr.zap = getZap(loggr.app, !loggr.disableStdout)
		loggr.writer = loggr.zapWrite
	} else {
		loggr.zero = getZero(loggr.app, !loggr.disableStdout)
		loggr.writer = loggr.zeroWrite
	}
	if err := loggr.finalizeEssentials(); err != nil {
		panic("unable to initialize logger")
	}
	return loggr
}

func MustGet(opts ...Option) IDotpeLogger {
	compress, _ = strconv.ParseBool(os.Getenv("LOGCOMPRESS"))
	loggeronce.Do(func() {
		logger = getlogger(opts...)
	})
	return logger
}

func DefaultLogWriter() ILogWriter { return defaultLogger }

func LogWriter() ILogWriter { return logger }
