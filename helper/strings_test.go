package helper

import "testing"

func TestStrJoin(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		{name: "t1", args: args{s: []string{"a", "b", "c"}}, wantStr: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := StrJoin(tt.args.s...); gotStr != tt.wantStr {
				t.Errorf("StrJoin() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}
