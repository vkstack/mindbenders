package logging

import (
	"fmt"
	"testing"

	"gitlab.com/dotpe/mindbenders/corel"
)

func Test_dlogger_Write(t *testing.T) {
	logger := MustGet(WithAppInfo("testing")).(*dlogger)
	fmt.Println(logger)
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
				cb:         InfoLevel,
				MessageKey: "testing",
			},
		},
	}
	ctx := corel.NewCorelCtx("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.WriteLogs(ctx, tt.args.fields, tt.args.cb, tt.args.MessageKey)
			// logger.zapWrite(tt.args.fields, tt.args.cb, tt.args.MessageKey)
		})
	}
}
