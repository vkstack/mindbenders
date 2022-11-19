package logging

import "github.com/sirupsen/logrus"

type IHookContainer interface {
	GetHook() (logrus.Hook, error)
	MustGetHook() logrus.Hook
}
