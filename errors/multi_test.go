package errors

import (
	"reflect"
	"testing"
)

var e4 error

func TestDefaultMultiError(t *testing.T) {
	tests := []struct {
		name string
		want MultiError
	}{
		{
			name: "test-1",
			want: &multierror{},
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
			e := &multierror{}
			e.AddErrors(tt.fields.errs...)
			if got := e.IsNil(); got != tt.want {
				t.Errorf("MultiError.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMultiError(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name string
		args args
		want MultiError
	}{
		{
			name: "test-1",
			args: args{[]error{e1, e2, e3}},
			want: &multierror{e1, e2, e3},
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

func Test_multierror_AddErrors(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name string
		e    *multierror
		args args
		want MultiError
	}{
		{
			name: "test-1",
			args: args{[]error{e1, e2, e3, e4, e0, e4}},
			want: &multierror{e1, e2, e3, e0},
		},
		{
			name: "test-2",
			args: args{[]error{e1, e2, e3, e4, e0, e4}},
			want: &multierror{e1, e2, e3, e0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e = &multierror{}
			if got := tt.e.AddErrors(tt.args.errs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("multierror.AddErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
