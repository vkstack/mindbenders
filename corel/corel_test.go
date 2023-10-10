package corel

import (
	"context"
	"reflect"
	"testing"
)

func TestNewCorelCtx(t *testing.T) {
	type args struct {
		sessionId string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "test-1",
			args: args{
				sessionId: "test-session",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrphenContext(tt.args.sessionId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCorelCtx() = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}
