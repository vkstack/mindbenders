package errors

import (
	"fmt"
	"testing"
)

func TestNewMultiError(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				errs: []error{nil, nil},
			},
			wantErr: false,
		},
		{
			name: "test-2",
			args: args{
				errs: []error{nil, nil, New("new testing errors")},
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewMultiError(tt.args.errs...)
			fmt.Println(err, err == nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMultiError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
