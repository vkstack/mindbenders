package errors

import (
	"testing"
)

var (
	e0 = New("base-error")
	e1 = WrapMessage(e0, "wrapped e0")
	e2 = WrapMessage(e1, "wrapped e1")
	e3 = WrapMessage(e2, "wrapped e2")
)

func Test_base_Error(t *testing.T) {
	tests := []struct {
		name  string
		error error
		want  string
	}{
		{
			name:  "test-0",
			error: e0,
			want:  "base-error",
		},
		{
			name:  "test-1",
			error: e1,
			want:  "wrapped e0",
		},
		{
			name:  "test-2",
			error: e2,
			want:  "wrapped e1",
		},
		{
			name:  "test-3",
			error: e3,
			want:  "wrapped e2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if e, ok := tt.error.(BaseError); !ok {
				return
			} else {
				if got := e.Error(); got != tt.want {
					t.Errorf("base.Error() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_base_String(t *testing.T) {
	tests := []struct {
		name  string
		error error
		want  string
	}{
		{
			name:  "test-0",
			error: e0,
			want:  "base-error",
		},
		{
			name:  "test-1",
			error: e1,
			want:  "wrapped e0\nbase-error",
		},
		{
			name:  "test-2",
			error: e2,
			want:  "wrapped e1\nwrapped e0\nbase-error",
		},
		{
			name:  "test-3",
			error: e3,
			want:  "wrapped e2\nwrapped e1\nwrapped e0\nbase-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if e, ok := tt.error.(BaseError); !ok {
				return
			} else {
				var got string
				if got = e.String(); got != tt.want {
					t.Errorf("base.String() = %v, want %v", got, tt.want)
				} else {
					t.Logf(got)
				}
			}
		})
	}
}

func TestCause(t *testing.T) {
	tests := []struct {
		name string
		error,
		want error
	}{
		{
			name:  "test-0",
			error: e0,
			want:  e0,
		},
		{
			name:  "test-1",
			error: e1,
			want:  e0,
		},
		{
			name:  "test-2",
			error: e2,
			want:  e0,
		},
		{
			name:  "test-3",
			error: e3,
			want:  e0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cause(tt.error); got != tt.want {
				t.Errorf("base.Cause() error = %v, wantErr %v", got, tt.want)
			} else {
				t.Logf("%v", got)
			}
		})
	}
}

func TestUnWrap(t *testing.T) {
	tests := []struct {
		name string
		error,
		want error
	}{
		{
			name:  "test-0",
			error: e0,
			want:  nil,
		},
		{
			name:  "test-1",
			error: e1,
			want:  e0,
		},
		{
			name:  "test-2",
			error: e2,
			want:  e1,
		},
		{
			name:  "test-3",
			error: e3,
			want:  e2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnWrap(tt.error); got != tt.want {
				t.Errorf("base.Cause() error = %v, wantErr %v", got, tt.want)
			} else {
				t.Logf("%v", got)
			}
		})
	}
}
