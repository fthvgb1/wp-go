package strings

import (
	"golang.org/x/exp/constraints"
	"testing"
)

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
			if gotStr := Join(tt.args.s...); gotStr != tt.wantStr {
				t.Errorf("Join() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestToInteger(t *testing.T) {
	type args[T constraints.Integer] struct {
		s string
		z T
	}
	type testCase[T constraints.Integer] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int64]{
		{
			name: "t1",
			args: args[int64]{
				"10",
				0,
			},
			want: int64(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInteger[int64](tt.args.s, tt.args.z); got != tt.want {
				t.Errorf("StrToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
