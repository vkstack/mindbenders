package logging

import (
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMustGetFileHook(t *testing.T) {
	type args struct {
		app string
	}
	tests := []struct {
		name string
		args args
		want logrus.Hook
	}{
		{
			name: "test-0",
			args: args{
				"comm",
			},
		},
		// TODO: Add test cases.
	}
	os.Setenv("LOGCOMPRESS", "true")
	os.Setenv("LOGSIZE", "1000")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustGetFileHook(tt.args.app)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustGetFileHook() = %v, want %v", got, tt.want)
			}
		})
	}
}
