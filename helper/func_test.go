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
