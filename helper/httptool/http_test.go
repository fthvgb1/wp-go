package httptool

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	type args struct {
		u       string
		q       map[string]any
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
				q: map[string]any{
					"p":                    "2",
					"XDEBUG_SESSION_START": "34343",
					"a":                    []int{2, 3},
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
			gotR, gotCode, err := GetString(tt.args.u, tt.args.q, tt.args.a...)
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

func TestPostWwwString(t *testing.T) {
	type args struct {
		u       string
		form    map[string]any
		timeout int64
		a       []any
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				u: "http://wp.test?XDEBUG_SESSION_START=34244",
				form: map[string]any{
					"aa":   "bb",
					"bb[]": []int{1, 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, _, err := PostWwwString(tt.args.u, tt.args.form, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Post() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		u       string
		types   int
		form    map[string]any
		timeout int64
		a       []any
	}
	tests := []struct {
		name    string
		args    args
		wantRes *http.Response
		wantErr bool
	}{
		{
			name: "form-data",
			args: args{
				u:     "http://wp.test?XDEBUG_SESSION_START=3424",
				types: 3,
				form: map[string]any{
					"ff": "xxxff",
				},
				timeout: 0,
				a:       nil,
			},
		},
		{
			name: "raw-json",
			args: args{
				u:     "http://wp.test?XDEBUG_SESSION_START=3424",
				types: 3,
				form: map[string]any{
					"ff": "xxxff",
					"kk": 1,
				},
				timeout: 0,
				a:       nil,
			},
		},
		{
			name: "binary",
			args: args{
				u:     "http://wp.test?XDEBUG_SESSION_START=3424",
				types: 4,
				form: map[string]any{
					"binary": []byte("ssssskkkkkk"),
				},
				a: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := Post(tt.args.u, tt.args.types, tt.args.form, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Post() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

type res struct {
	Code    int    `json:"Code,omitempty"`
	Message string `json:"Message" json:"Message,omitempty"`
}

func TestPostToJsonAny(t *testing.T) {
	type args struct {
		u     string
		types int
		form  map[string]any
		a     []any
	}
	type testCase[T any] struct {
		name     string
		args     args
		wantR    T
		wantCode int
		wantErr  bool
	}

	tests := []testCase[res]{
		{
			name: "res",
			args: args{
				u:     "http://wp.test?XDEBUG_SESSION_START=3424",
				types: 1,
				a:     []any{3 * time.Second, map[string]string{"user-agent": "httptool"}},
			},
			wantR: res{
				200, "ok",
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotCode, err := PostToJsonAny[res](tt.args.u, tt.args.types, tt.args.form, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostToJsonAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("PostToJsonAny() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotCode != tt.wantCode {
				t.Errorf("PostToJsonAny() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestGetToJsonAny(t *testing.T) {
	type args struct {
		u string
		q map[string]any
		a []any
	}
	type testCase[T any] struct {
		name     string
		args     args
		wantR    T
		wantCode int
		wantErr  bool
	}
	tests := []testCase[res]{
		{
			name: "t1",
			args: args{
				u: "http://wp.test?XDEBUG_SESSION_START=3424",
				q: map[string]any{
					"jjj": "ssss",
					"fff": []int{1, 2, 3},
				},
			},
			wantR: res{
				200, "ok",
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotCode, err := GetToJsonAny[res](tt.args.u, tt.args.q, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToJsonAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetToJsonAny() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotCode != tt.wantCode {
				t.Errorf("GetToJsonAny() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
