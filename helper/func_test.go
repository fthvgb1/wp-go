package helper

import (
	"context"
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

func BenchmarkStructColumnToSlice(b *testing.B) {
	y := getX()
	fmt.Println(y)
	b.ResetTimer()
	//b.N = 2
	for i := 0; i < 1; i++ {
		StructColumnToSlice[int, *x](y, "Id")
	}
}

func TestStructColumnToSlice(t *testing.T) {
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
			if gotR := StructColumnToSlice[uint64, x](tt.args.arr, tt.args.field); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("StructColumnToSlice() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestToAny(t *testing.T) {
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

func TestCutUrlHost(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "http",
			args: args{"http://xx.yy/xxoo?ss=fff"},
			want: "/xxoo?ss=fff",
		}, {
			name: "https",
			args: args{"https://xx.yy/xxoo?ff=fff"},
			want: "/xxoo?ff=fff",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CutUrlHost(tt.args.u); got != tt.want {
				t.Errorf("CutUrlHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	v := 0
	want := 1
	t.Run("int", func(t *testing.T) {
		if got := Defaults(v, want); !reflect.DeepEqual(got, want) {
			t.Errorf("Defaults() = %v, want %v", got, want)
		}
	})
	{
		v := ""
		want := "a"
		t.Run("string", func(t *testing.T) {
			if got := Defaults(v, want); !reflect.DeepEqual(got, want) {
				t.Errorf("Defaults() = %v, want %v", got, want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	{
		name := "bool"
		args := true
		want := true
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "int"
		args := 0
		want := false
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "int"
		args := 1
		want := true
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "string"
		args := "1"
		want := true
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "string"
		args := "0"
		want := false
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "string"
		args := ""
		want := false
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
	{
		name := "float"
		args := 0.2
		want := true
		t.Run(name, func(t *testing.T) {
			if got := ToBool(args); got != want {
				t.Errorf("ToBool() = %v, want %v", got, want)
			}
		})
	}
}

func TestIsZeros(t *testing.T) {
	tt := struct {
		name string
		args struct {
			v struct{ a string }
		}
		want bool
	}{
		name: "t1",
		args: struct{ v struct{ a string } }{v: struct{ a string }{a: ""}},
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := IsZeros(tt.args.v); got != tt.want {
			t.Errorf("IsZeros() = %v, want %v", got, tt.want)
		}
	})
}

func TestGetValFromContext(t *testing.T) {
	type args[K any, V any] struct {
		ctx      context.Context
		k        K
		defaults V
	}
	type testCase[K any, V any] struct {
		name string
		args args[K, V]
		want V
	}
	tests := []testCase[string, int]{
		{
			name: "t1",
			args: args[string, int]{
				ctx:      context.WithValue(context.Background(), "kk", 1),
				k:        "kk",
				defaults: 0,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContextVal(tt.args.ctx, tt.args.k, tt.args.defaults); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContextVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAnyVal(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := "string"
		if got := GetAnyVal(any("string"), "s"); !reflect.DeepEqual(got, want) {
			t.Errorf("GetAnyVal() = %v, want %v", got, want)
		}
		want = "s"
		if got := GetAnyVal(any(1), "s"); !reflect.DeepEqual(got, want) {
			t.Errorf("GetAnyVal() = %v, want %v", got, want)
		}
	})
}
