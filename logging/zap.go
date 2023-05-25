package logging

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func (l *dlogger) enzap(fields logrus.Fields) (zfields []zap.Field) {
	for k, v := range fields {
		zfields = append(zfields, zap.Any(k, v))
	}
	return
}
