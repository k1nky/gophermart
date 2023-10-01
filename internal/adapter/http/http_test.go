package http

import (
	"net/http"
	"testing"
)

func TestAdapter_Register(t *testing.T) {
	type fields struct {
		auth AuthService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adapter{
				auth: tt.fields.auth,
			}
			a.Register(tt.args.w, tt.args.r)
		})
	}
}
