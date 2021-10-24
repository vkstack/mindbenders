package logging

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type Option func(dlogger *dlogger)

func WithAppInfo(app, env, wd string) Option {
	return func(dlogger *dlogger) {
		dlogger.app = app
		dlogger.env = env
		dlogger.wd = wd
	}
}

func WithHookContainer(hookContainer ILogConfig) Option {
	hook, err := hookContainer.getHook()
	if err != nil {
		return nil
	}
	return WithHook(hook)
}

func WithHook(hook logrus.Hook) Option {
	return func(dlogger *dlogger) {
		dlogger.logger = logrus.New()
		dlogger.logger.SetNoLock()
		dlogger.logger.Hooks.Add(hook)
		if dlogger.env != "dev" {
			dlogger.logger.Out = ioutil.Discard
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
