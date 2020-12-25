package interfaces

import (
	"context"

	"github.com/sirupsen/logrus"
)

//ILogger ...
type ILogger interface {
	WriteLogs(context.Context, logrus.Fields, logrus.Level, string, ...interface{})
}
