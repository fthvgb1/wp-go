package httptool

import (
	"testing"
)

func TestGetString(t *testing.T) {
	type args struct {
		u       string
		q       map[string]string
		timeout int64
		a       []any
	}
	tests := []struct {
		name     string
		args     args
		wantR    string
		wantCode int
		wantErr  bool
	}{
		{
			name: "wp.test",
			args: args{
				u: "http://wp.test",
				q: map[string]string{
					"p":                    "2",
					"XDEBUG_SESSION_START": "34343",
				},
				timeout: 3,
			},
			wantR:    `{"XDEBUG_SESSION_START":"34343","p":"2"}`,
			wantCode: 200,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotCode, err := GetString(tt.args.u, tt.args.q, tt.args.timeout, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("GetString() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotCode != tt.wantCode {
				t.Errorf("GetString() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
