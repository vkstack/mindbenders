package errors

import (
	"reflect"
	"testing"
)

var e4 error

func TestNewMultiError(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name string
		args args
		want *MultiError
	}{
		{
			name: "test-1",
			args: args{
				errs: []error{e1, e2, e3},
			},
			want: &MultiError{errs: []error{e1, e2, e3}},
		},
		{
			name: "test-2",
			args: args{
				errs: []error{e1, e2, e3, e4},
			},
			want: &MultiError{errs: []error{e1, e2, e3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMultiError(tt.args.errs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMultiError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultMultiError(t *testing.T) {
	tests := []struct {
		name string
		want *MultiError
	}{
		{
			name: "test-1",
			want: &MultiError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultMultiError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultMultiError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiError_AddErrors(t *testing.T) {
	type fields struct {
		errs []error
	}
	type args struct {
		errs []error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MultiError
	}{
		{
			name:   "test-1",
			fields: fields{[]error{e4, e1}},
			args:   args{[]error{e2, e3, e4, e0, e4}},
			want:   &MultiError{[]error{e1, e2, e3, e0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewMultiError(tt.fields.errs...)
			if got := e.AddErrors(tt.args.errs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultiError.AddErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiError_IsNil(t *testing.T) {
	type fields struct {
		errs []error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test-1",
			fields: fields{[]error{e4, e4, e4}},
			want:   true,
		},
		{
			name:   "test-2",
			fields: fields{[]error{e4, e1, e4}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewMultiError(tt.fields.errs...)
			if got := e.IsNil(); got != tt.want {
				t.Errorf("MultiError.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}
