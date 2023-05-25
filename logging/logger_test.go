package logging

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

func Test_dlogger_Write(t *testing.T) {
	logger := MustGet(WithAppInfo("testing")).(*dlogger)
	fmt.Println(logger)
	type args struct {
		fields     logrus.Fields
		cb         zapcore.Level
		MessageKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				fields:     logrus.Fields{"test-field-1": "test-field-value-1"},
				cb:         zapcore.InfoLevel,
				MessageKey: "testing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Write(tt.args.fields, zapcore.Level(tt.args.cb), tt.args.MessageKey)
		})
	}
}
