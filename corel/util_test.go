package corel

import (
	"context"
	"reflect"
	"testing"
)

func TestEncodeCorel(t *testing.T) {
	type args struct {
		corelId *CoRelationId
	}
	corelId := &CoRelationId{SessionId: "sessionId"}
	corelId.init(context.Background())
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{
				corelId: corelId.child(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeCorel(tt.args.corelId)
			x := DecodeCorelationId(got)
			ok := reflect.DeepEqual(&x, tt.args.corelId)
			if ok && got != tt.want {
				t.Errorf("EncodeCorel() = %v, want %v", got, tt.want)
			}
		})
	}
}
