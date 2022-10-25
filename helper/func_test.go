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
		{
			name: "t4",
			args: args{
				start: 0,
				end:   -5,
				step:  -1,
			},
			want: []int{0, -1, -2, -3, -4, -5},
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

func TestCloseHtmlTag(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{str: `<pre class="wp-block-preformatted">GRANT privileges ON databasename.tablename TO 'username'@'h...<p class="read-more"><a href="/p/305">继续阅读</a></p>`},
			want: "</pre>",
		},
		{
			name: "t2",
			args: args{str: `<pre><div>`},
			want: "</div></pre>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloseHtmlTag(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloseHtmlTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clearTag(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{s: []string{"<pre>", "<p>", "<span>", "</span>"}},
			want: []string{"<pre>", "<p>"},
		},
		{
			name: "t2",
			args: args{s: []string{"<pre>", "</pre>", "<div>", "<span>", "</span>"}},
			want: []string{"<div>"},
		},
		{
			name: "t3",
			args: args{s: []string{"<pre>", "</pre>"}},
			want: []string{},
		},
		{
			name: "t4",
			args: args{s: []string{"<pre>", "<p>"}},
			want: []string{"<pre>", "<p>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClearClosedTag(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClearClosedTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int, int) int
		r   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "t1",
			args: args{arr: RangeSlice(1, 10, 1), fn: func(i int, i2 int) int {
				return i + i2
			}},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceReduce(tt.args.arr, tt.args.fn, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceReduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceFilter(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{arr: RangeSlice(1, 10, 1), fn: func(i int) bool {
				return i > 4
			}},
			want: RangeSlice(5, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceFilter(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceMap(t *testing.T) {
	type args struct {
		arr []int8
		fn  func(int8) int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr: RangeSlice[int8](1, 10, 1),
				fn: func(i int8) int {
					return int(i)
				},
			},
			want: RangeSlice(1, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceMap(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReverse(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{arr: RangeSlice(1, 10, 1)},
			want: RangeSlice(10, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceReverse(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInterface(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "t1",
			args: args{v: 1},
			want: any(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAny(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceSelfReverse(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "t1",
			args: args{
				arr: RangeSlice(1, 10, 1),
			},
			want: RangeSlice(10, 1, -1),
		}, {
			name: "t2",
			args: args{
				arr: RangeSlice(1, 9, 1),
			},
			want: RangeSlice(9, 1, -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceSelfReverse(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceSelfReverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceToMap(t *testing.T) {
	type ss struct {
		id int
		v  string
	}
	type args struct {
		arr         []ss
		fn          func(ss) (int, ss)
		isCoverPrev bool
	}
	tests := []struct {
		name string
		args args
		want map[int]ss
	}{
		{
			name: "t1",
			args: args{
				arr: []ss{{1, "k1"}, {2, "v2"}, {2, "v3"}},
				fn: func(s ss) (int, ss) {
					return s.id, s
				},
				isCoverPrev: true,
			},
			want: map[int]ss{1: {1, "k1"}, 2: {2, "v3"}},
		}, {
			name: "t2",
			args: args{
				arr: []ss{{1, "k1"}, {2, "v2"}, {2, "v3"}},
				fn: func(s ss) (int, ss) {
					return s.id, s
				},
				isCoverPrev: false,
			},
			want: map[int]ss{1: {1, "k1"}, 2: {2, "v2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceToMap(tt.args.arr, tt.args.fn, tt.args.isCoverPrev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleSliceToMap(t *testing.T) {
	type args struct {
		arr []int
		fn  func(int) int
	}
	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "t1",
			args: args{arr: []int{1, 2, 3}, fn: func(i int) int {
				return i
			}},
			want: map[int]int{1: 1, 2: 2, 3: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleSliceToMap(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleSliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandNum(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "t1",
			args: args{
				start: 1,
				end:   2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := RandNum(tt.args.start, tt.args.end)
				if got > tt.args.end || got < tt.args.start {
					t.Errorf("RandNum() = %v, range error", got)
				}
			}
		})
	}
}

func TestSampleSort(t *testing.T) {
	type args struct {
		arr []int
		fn  func(i, j int) bool
	}
	tests := []struct {
		name  string
		args  args
		wantR []int
	}{
		{
			name: "t1",
			args: args{
				arr: []int{3, 5, 6, 1},
				fn: func(i, j int) bool {
					return i < j
				},
			},
			wantR: []int{1, 3, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if SampleSort(tt.args.arr, tt.args.fn); !reflect.DeepEqual(tt.args.arr, tt.wantR) {
				t.Errorf("SampleSort() = %v, want %v", tt.args.arr, tt.wantR)
			}
		})
	}
}
