package corel

import (
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
			fields: fields{Auth: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImVscGl6emFyaW9fMWsiLCJUZW5hbnRJZCI6NjgxNiwiU3RvcmVJZCI6MTM2MTIsIkZlYXR1cmVSb2xlSWQiOiJzdG9yZV9vdmVydmlldyxzdG9yZV90cmFuc2FjdGlvbnMsc3RvcmVfaXRlbV9zYWxlc19yZWdpc3RlcixtZXJjaGFudF9jYXRhbG9nX3N0b3JlLHN0b3JlX2NhdGFsb2csa2l0Y2hlbl9sb2dvdXQsbWVyY2hhbnRfY2F0YWxvZ19hZGRfaXRlbSxtZXJjaGFudF9jYXRhbG9nX3VwZGF0ZV9pdGVtLGRlbGV0ZV9pdGVtLGRlbGV0ZV9jYXRlZ29yeSxtZXJjaGFudF9jYXRhbG9nLHJlamVjdF9vcmRlcixraXRjaGVuX29yZGVyX3ByaW50LG9yZGVyX2FsbF9zdGF0ZXNfdmlldyxuZXdfb3JkZXJfcmVqZWN0aW9uLGFjY2VwdF9vcmRlcl9yZWplY3Rpb24iLCJFeHBpcnlUaW1lIjoiMjAyMS0wNy0yNiAwNTowMDowMCIsIklzc3VlVGltZSI6IjIwMjEtMDctMjUgMTI6MjY6MTgiLCJVc2VyVHlwZSI6ImtpdGNoZW4iLCJleHAiOjE2MjcyNTU4MDAsImlzcyI6ImRvdHBlS2l0Y2hlbiJ9.F9-6a_pglutP9mUVHmIwn7pFiEiF-cfJrPm9BenTrdk"},
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
