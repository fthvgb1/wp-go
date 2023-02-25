package safety

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/taskPools"
	"reflect"
	"testing"
)

func TestMap_Load(t *testing.T) {
	type args[K comparable] struct {
		key K
	}
	m := NewMap[int, int]()
	m.Store(1, 1)
	type testCase[K comparable, V any] struct {
		name      string
		m         *Map[K, V]
		args      args[K]
		wantValue V
		wantOk    bool
	}
	tests := []testCase[int, int]{
		{
			name:      "t1",
			m:         m,
			args:      args[int]{1},
			wantValue: 1,
			wantOk:    true,
		},
	}
	p := taskPools.NewPools(10)
	var a0, a1 []int
	for i := 0; i < 15000; i++ {
		v := number.Rand(0, 2)
		if 1 == v {
			a1 = append(a1, 1)
		} else if 0 == v {
			m.Flush()
		} else {
			a0 = append(a0, 0)
		}
		p.Execute(func() {
			if 1 == v {
				m.Load(number.Rand(2, 1000))
			} else {
				m.Store(number.Rand(2, 1000), number.Rand(1, 1000))
			}
		})

	}
	fmt.Println(len(a0), len(a1), m.Len())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
