package errors

import "testing"

func Test_base_Error(t *testing.T) {
	type fields struct {
		msg   string
		cause error
		code  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"test-1",
			fields{
				msg: "test-1",
			},
			"test-1",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &base{
				msg:   tt.fields.msg,
				cause: tt.fields.cause,
				code:  tt.fields.code,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("base.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
