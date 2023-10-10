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

func Test_getJWTSession(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{
				token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJhaW5ib3cxIiwiVGVuYW50SWQiOjU2NjAsIkZlYXR1cmVSb2xlSWQiOiJzdG9yZV9jYXRhbG9nLG1lcmNoYW50X2NhdGFsb2dfaXRlbV9zdGF0dXMsbWVyY2hhbnRfY2F0YWxvZyxzdG9yZV9oZWFsdGgsbWVyY2hhbnRfaXRlbV9zYWxlc19yZWdpc3RlcixmaW5hbmNlX3NldHRsZW1lbnRfZ2V0LHN0b3Jlc19kZXRhaWxzLG1lcmNoYW50X3N0b3JlX2ZsYWcscXJjb2RlX2dlbixzdG9yZV90aW1pbmcsbWVyY2hhbnRfY2F0YWxvZ19vcHRfZGVsZXRlLG1lcmNoYW50X2NhdGFsb2dfb3B0X3N0YXR1cyxtZXJjaGFudF9jYXRhbG9nX3Zhcl9zdGF0dXMsbWVyY2hhbnRfY2F0YWxvZ192YXJfZGVsZXRlLGJpbGxfYWRqdXN0bWVudCx1cGxvYWRfaXRlbV9pbWFnZSxkZWxldGVfY2F0ZWdvcnksbWVyY2hhbnRfY2F0YWxvZ19hZGRfaXRlbSxtZXJjaGFudF9jYXRhbG9nX3VwZGF0ZV9pdGVtLGJ1bGtfdXBsb2FkX211bHRpX3N0b3JlLHJlc2V0X3Bhc3N3b3JkLGRlbGV0ZV9pdGVtLG1lcmNoYW50X3RyYW5zYWN0aW9ucyxtZXJjaGFudF9vdmVydmlldyxtZXJjaGFudF9jdXN0b21lcl9leHBlcmllbmNlLG1lcmNoYW50X2hlYWx0aCx1cGRhdGVfY2F0ZWdvcnksbWVyY2hhbnRfbWFya2V0aW5nLHdhX21lcmNoYW50X21hcmtldGluZyIsIlVzZXJUeXBlIjoibWVyY2hhbnQiLCJVc2VySUQiOjUwMDE1Nywic2Vzc2lvbklEIjoiRUk0TVVDRFhlN0VnMDhUdCIsIlN0b3JlSUQiOjAsImRldmljZUlEIjoiIiwiSXNBZG1pbiI6ZmFsc2UsIlBob25lIjoiIiwiUGFydG5lckFwcCI6ZmFsc2UsIlBob25lTG9naW4iOmZhbHNlLCJhcHBOYW1lIjoiIiwiZXhwIjoxNjk2OTk0MzI5LCJpc3MiOiJkb3RwZSJ9.G8jTSqq12GJsMBud9SdIoskG80Fo8-jxgU1W4--bDus",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getJWTSession(tt.args.token)
			if got != tt.want {
				t.Errorf("getJWTSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
