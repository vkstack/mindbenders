package logging

import (
	"os"
)

type Option func(dlogger *dlogger)

func WithAppInfo(app string) Option {
	return func(dlogger *dlogger) {
		dlogger.app = app
		dlogger.env = os.Getenv("ENV")
	}
}

func WithAccessLogOptions(opts ...accessLogOption) Option {
	return func(dlogger *dlogger) {
		dlogger.accopts = append(dlogger.accopts, accessLogOptionBasic(dlogger.app))
		dlogger.accopts = append(dlogger.accopts, opts...)
	}
}

func WithLogOptions(opts ...logOption) Option {
	return func(dlogger *dlogger) {
		dlogger.loptions = append(dlogger.loptions, logOptionBasic)
		dlogger.loptions = append(dlogger.loptions, opts...)
	}
}

func WithZap(dlogger *dlogger) { dlogger.iszap = true }

func WithZero(dlogger *dlogger) { dlogger.iszap = false }

func DisabledStdLogging(dlogger *dlogger) { dlogger.disableStdout = true }
