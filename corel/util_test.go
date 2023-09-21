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
	corelId := NewCorelId("sessionId")
	corelId.init(context.Background())
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{
				corelId: corelId.Child(),
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

func Test_corel(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *CoRelationId
		wantErr bool
	}{
		{
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCorelId(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("corel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("corel() = %v, want %v", got, tt.want)
			}
		})
	}
}
