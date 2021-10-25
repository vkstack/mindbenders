package logging

import "github.com/sirupsen/logrus"

type ILogConfig interface {
	getHook() (logrus.Hook, error)
}
