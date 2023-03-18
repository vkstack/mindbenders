package logging

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type Option func(dlogger *dlogger)

func WithAppInfo(app string) Option {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return func(dlogger *dlogger) {
		dlogger.app = app
		dlogger.env = os.Getenv("ENV")
		dlogger.wd = wd
	}
}

func WithHook(hook logrus.Hook) Option {
	if hook == nil {
		return nil
	}
	return func(dlogger *dlogger) {
		dlogger.logger = logrus.New()
		dlogger.logger.SetNoLock()
		dlogger.logger.Hooks.Add(hook)
		if dlogger.env != "dev" {
			dlogger.logger.Out = io.Discard
		}
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
