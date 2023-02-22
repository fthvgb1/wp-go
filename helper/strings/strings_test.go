package strings

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
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

func TestBuilder_WriteString(t *testing.T) {
	type fields struct {
		Builder *strings.Builder
	}
	type args struct {
		s []string
	}
	//s :=NewBuilder()
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "t1",
			fields: fields{
				Builder: &strings.Builder{},
			},
			args:      args{s: []string{"11", "22", "不"}},
			wantErr:   false,
			wantCount: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				Builder: tt.fields.Builder,
			}
			gotCount := b.WriteString(tt.args.s...)

			if gotCount != tt.wantCount {
				t.Errorf("WriteString() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
			fmt.Println(b.String())
		})
	}
}

func BenchmarkBuilder_SprintfXX(b *testing.B) {
	s := NewBuilder()
	for i := 0; i < b.N; i++ {
		s.Sprintf("%s %s %s", "a", "b", "c")
		_ = s.String()
	}
}

func BenchmarkSPrintfXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s %s %s", "a", "b", "c")
	}
}

func BenchmarkStrJoinXX(b *testing.B) {
	s := strings.Builder{}
	for i := 0; i < b.N; i++ {
		s.WriteString("a ")
		s.WriteString("b ")
		s.WriteString("c ")
		_ = s.String()
	}
}
func BenchmarkBuilderJoinXX(b *testing.B) {
	s := NewBuilder()
	for i := 0; i < b.N; i++ {
		s.WriteString("a ", "b ", "c")
		_ = s.String()
	}
}
