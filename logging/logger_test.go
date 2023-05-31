package logging

import (
	"testing"

	"gitlab.com/dotpe/mindbenders/corel"
)

func Test_dlogger_Write(t *testing.T) {
	// os.Setenv("LOGLEVEL", "debug")
	logger := getlogger(WithAppInfo("testing"), WithZero)
	logger1 := getlogger(WithAppInfo("testing"), WithZap)
	type args struct {
		fields     Fields
		cb         Level
		MessageKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				fields:     Fields{"test-field-1": "test-field-value-1"},
				cb:         FatalLevel,
				MessageKey: "testing",
			},
		},
	}
	for _, tt := range tests {
		ctx := corel.NewCorelCtx("testing")
		t.Run(tt.name, func(t *testing.T) {
			logger.Debug(ctx, tt.args.fields, tt.args.MessageKey)
			logger.Info(ctx, tt.args.fields, tt.args.MessageKey)
			logger.Warn(ctx, tt.args.fields, tt.args.MessageKey)
			logger.Error(ctx, tt.args.fields, tt.args.MessageKey)
			logger.Panic(ctx, tt.args.fields, tt.args.MessageKey)
			logger.Fatal(ctx, tt.args.fields, tt.args.MessageKey)

			logger1.Log(ctx, tt.args.fields, FatalLevel, tt.args.MessageKey)
			logger.Log(ctx, tt.args.fields, FatalLevel, tt.args.MessageKey)
			logger1.Log(ctx, tt.args.fields, PanicLevel, tt.args.MessageKey)
			logger.Log(ctx, tt.args.fields, PanicLevel, tt.args.MessageKey)

			logger1.Log(ctx, tt.args.fields, InfoLevel, tt.args.MessageKey)
			logger.Log(ctx, tt.args.fields, InfoLevel, tt.args.MessageKey)
		})
	}
}
