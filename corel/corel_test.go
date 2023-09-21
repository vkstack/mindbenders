package corel

import (
	"context"
	"reflect"
	"testing"
)

func TestCoRelationId_loadAuth(t *testing.T) {
	type fields struct {
		Auth string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "test",
			fields: fields{Auth: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkJhcmlzdGFfa2l0Y2hlbiIsIlRlbmFudElkIjo1LCJTdG9yZUlkIjo5LCJGZWF0dXJlUm9sZUlkIjoic3RvcmVfb3ZlcnZpZXcsc3RvcmVfdHJhbnNhY3Rpb25zLHN0b3JlX2l0ZW1fc2FsZXNfcmVnaXN0ZXIsbWVyY2hhbnRfY2F0YWxvZ19zdG9yZSxzdG9yZV9jYXRhbG9nLGtpdGNoZW5fbG9nb3V0LG1lcmNoYW50X2NhdGFsb2dfYWRkX2l0ZW0sbWVyY2hhbnRfY2F0YWxvZ191cGRhdGVfaXRlbSxkZWxldGVfaXRlbSxkZWxldGVfY2F0ZWdvcnksbWVyY2hhbnRfY2F0YWxvZyxyZWplY3Rfb3JkZXIsa2l0Y2hlbl9vcmRlcl9wcmludCxvcmRlcl9hbGxfc3RhdGVzX3ZpZXcsbmV3X29yZGVyX3JlamVjdGlvbixhY2NlcHRfb3JkZXJfcmVqZWN0aW9uLGNvbXBsZXRlX29yZGVyX3JlamVjdGlvbiIsIkV4cGlyeVRpbWUiOiIyMDIxLTA4LTIwIDA1OjAwOjAwIiwiSXNzdWVUaW1lIjoiMjAyMS0wOC0xOSAyMzoxNzoxOCIsIlVzZXJUeXBlIjoia2l0Y2hlbiIsImV4cCI6MTYyOTQxNTgwMCwiaXNzIjoiZG90cGVLaXRjaGVuIn0.i-2C61Bwmopm1AQ3_BBxopzr73FSdEITXgSvzaNAo20"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corelid := &CoRelationId{
				Auth: tt.fields.Auth,
			}
			corelid.loadAuth()
		})
	}
}

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
