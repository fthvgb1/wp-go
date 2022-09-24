package helper

import (
	"fmt"
	"reflect"
	"testing"
)

type x struct {
	Id uint64
}

func c(x []*x) (r []uint64) {
	for i := 0; i < len(x); i++ {
		r = append(r, x[i].Id)
	}
	return
}

func getX() (r []*x) {
	for i := 0; i < 10; i++ {
		r = append(r, &x{
			uint64(i),
		})
	}
	return
}

func BenchmarkOr(b *testing.B) {
	y := getX()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c(y)
	}
}

func BenchmarkStructColumn(b *testing.B) {
	y := getX()
	fmt.Println(y)
	b.ResetTimer()
	//b.N = 2
	for i := 0; i < 1; i++ {
		StructColumn[int, *x](y, "Id")
	}
}

func TestStructColumn(t *testing.T) {
	type args struct {
		arr   []x
		field string
	}

	tests := []struct {
		name  string
		args  args
		wantR []uint64
	}{
		{name: "test1", args: args{
			arr: []x{
				{Id: 1},
				{2},
				{4},
				{6},
			},
			field: "Id",
		}, wantR: []uint64{1, 2, 4, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := StructColumn[uint64, x](tt.args.arr, tt.args.field); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("StructColumn() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestRangeSlice(t *testing.T) {
	type args struct {
		start int
		end   int
		step  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				start: 1,
				end:   5,
				step:  1,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "t2",
			args: args{
				start: 0,
				end:   5,
				step:  2,
			},
			want: []int{0, 2, 4},
		},
		{
			name: "t3",
			args: args{
				start: 1,
				end:   11,
				step:  3,
			},
			want: []int{1, 4, 7, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RangeSlice(tt.args.start, tt.args.end, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RangeSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestSlicePagination(t *testing.T) {
	arr := RangeSlice[int](1, 10, 1)
	type args struct {
		arr      []int
		page     int
		pageSize int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr:      arr,
				page:     1,
				pageSize: 2,
			},
			want: RangeSlice[int](1, 2, 1),
		}, {
			name: "t2",
			args: args{
				arr:      arr,
				page:     2,
				pageSize: 2,
			},
			want: RangeSlice[int](3, 4, 1),
		}, {
			name: "t3",
			args: args{
				arr:      arr,
				page:     4,
				pageSize: 3,
			},
			want: []int{10},
		}, {
			name: "t4",
			args: args{
				arr:      arr,
				page:     5,
				pageSize: 3,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlicePagination(tt.args.arr, tt.args.page, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlicePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	type args struct {
		str       string
		allowable string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				str:       "<p>ppppp<span>ffff</span></p><img />",
				allowable: "<p><img>",
			},
			want: "<p>pppppffff</p><img />",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripTags(tt.args.str, tt.args.allowable); got != tt.want {
				t.Errorf("StripTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripTagsX(t *testing.T) {
	type args struct {
		str       string
		allowable string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				str:       "<p>ppppp<span>ffff</span></p><img />",
				allowable: "<p><img>",
			},
			want: "<p>pppppffff</p><img />",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripTagsX(tt.args.str, tt.args.allowable); got != tt.want {
				t.Errorf("StripTagsX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkStripTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StripTags(`<p>ppppp<span>ffff</span></p><img />`, "<p><img>")
	}
}
func BenchmarkStripTagsX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StripTagsX(`<p>ppppp<span>ffff</span></p><img />`, "<p><img>")
	}
}
