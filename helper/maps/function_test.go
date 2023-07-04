package maps

import (
	"reflect"
	"testing"
)

func TestGetStrMapAnyVal(t *testing.T) {
	type args struct {
		key string
		v   map[string]any
	}
	type testCase[T any] struct {
		name  string
		args  args
		wantR T
		wantO bool
	}
	tests := []testCase[int]{
		{name: "t1", args: args{
			key: "k1",
			v: map[string]any{
				"k1": 1,
				"k2": 2,
			},
		}, wantR: 1, wantO: true},
		{name: "t2", args: args{
			key: "k2.kk",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: 10, wantO: true},
		{name: "t3", args: args{
			key: "k2.vv",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: 0, wantO: false},
		{name: "t4", args: args{
			key: "k3",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: 0, wantO: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotO := GetStrAnyVal[int](tt.args.v, tt.args.key)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetStrAnyVal() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotO != tt.wantO {
				t.Errorf("GetStrAnyVal() gotO = %v, want %v", gotO, tt.wantO)
			}
		})
	}
}

func TestGetStrMapAnyValWithAny(t *testing.T) {
	type args struct {
		key string
		v   map[string]any
	}
	tests := []struct {
		name  string
		args  args
		wantR any
		wantO bool
	}{
		{name: "t1", args: args{
			key: "k1",
			v: map[string]any{
				"k1": 1,
				"k2": 2,
			},
		}, wantR: any(1), wantO: true},
		{name: "t2", args: args{
			key: "k2.kk",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: any(10), wantO: true},
		{name: "t3", args: args{
			key: "k2.vv",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: nil, wantO: false},
		{name: "t4", args: args{
			key: "k3",
			v: map[string]any{
				"k1": 1,
				"k2": map[string]any{
					"kk": 10,
				},
			},
		}, wantR: nil, wantO: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotO := GetStrMapAnyValWithAny(tt.args.key, tt.args.v)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetStrMapAnyValWithAny() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotO != tt.wantO {
				t.Errorf("GetStrMapAnyValWithAny() gotO = %v, want %v", gotO, tt.wantO)
			}
		})
	}
}

func TestGetAnyAnyMapVal(t *testing.T) {
	m := map[any]any{
		1:   "1",
		"1": 2,
		"xxx": map[any]any{
			"sss": []int{1, 2, 3},
			5:     "5",
		},
	}

	t.Run("t1", func(t *testing.T) {
		wantR := []int{1, 2, 3}
		gotR, gotO := GetAnyAnyMapVal[[]int](m, "xxx", "sss")
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
		if gotO != true {
			t.Errorf("GetAnyAnyMapVal() gotO = %v, want %v", gotO, true)
		}
	})

	t.Run("t2", func(t *testing.T) {
		wantR := "5"
		gotR, gotO := GetAnyAnyMapVal[string](m, "xxx", 5)
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
		if gotO != true {
			t.Errorf("GetAnyAnyMapVal() gotO = %v, want %v", gotO, true)
		}
	})
}

func TestGetAnyAnyMapWithAny(t *testing.T) {
	m := map[any]any{
		1:   "1",
		"1": 2,
		"xxx": map[any]any{
			"sss": []int{1, 2, 3},
			5:     "5",
		},
	}

	t.Run("t1", func(t *testing.T) {
		wantR := any([]int{1, 2, 3})
		gotR, gotO := GetAnyAnyMapWithAny(m, "xxx", "sss")
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
		if gotO != true {
			t.Errorf("GetAnyAnyMapVal() gotO = %v, want %v", gotO, true)
		}
	})

	t.Run("t2", func(t *testing.T) {
		wantR := any("5")
		gotR, gotO := GetAnyAnyMapWithAny(m, "xxx", 5)
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
		if gotO != true {
			t.Errorf("GetAnyAnyMapVal() gotO = %v, want %v", gotO, true)
		}
	})
}

func TestGetAnyAnyValWithDefaults(t *testing.T) {
	m := map[any]any{
		1:   "1",
		"1": 2,
		"xxx": map[any]any{
			"sss": []int{1, 2, 3},
			5:     "5",
		},
	}

	t.Run("t1", func(t *testing.T) {
		wantR := []int{1, 2, 3}
		gotR := GetAnyAnyValWithDefaults[[]int](m, nil, "xxx", "sss")
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
	})

	t.Run("t2", func(t *testing.T) {
		wantR := "xxx"
		gotR := GetAnyAnyValWithDefaults[string](m, "xxx", "xxx", 55)
		if !reflect.DeepEqual(gotR, wantR) {
			t.Errorf("GetAnyAnyMapVal() gotR = %v, want %v", gotR, wantR)
		}
	})
}

func TestRecursiveSetAnyVal(t *testing.T) {
	type args[T any] struct {
		m map[any]any
		v T
		k []any
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[float64]{
		{
			name: "t1",
			args: args[float64]{
				m: map[any]any{},
				v: 3.4,
				k: []any{"m", 3, "key"},
			},
		},
	}
	target := map[any]any{"m": map[any]any{3: map[any]any{"key": 3.4}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAnyAnyVal(tt.args.m, tt.args.v, tt.args.k...)
			if !reflect.DeepEqual(tt.args.m, target) {
				t.Fatalf("not equal target")
			}
		})
	}
}

func TestRecursiveSetStrVal(t *testing.T) {
	type args[T any] struct {
		m map[string]any
		k string
		v T
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "t1",
			args: args[int]{
				m: map[string]any{},
				k: "aa.bb.cc",
				v: 1,
			},
		},
		{
			name: "t2",
			args: args[int]{
				m: map[string]any{"aa": map[string]any{}},
				k: "aa.bb.cc",
				v: 1,
			},
		},
		{
			name: "t3",
			args: args[int]{
				m: map[string]any{"aa": map[string]any{"bb": map[string]any{}}},
				k: "aa.bb.cc",
				v: 1,
			},
		},
	}
	target := map[string]any{"aa": map[string]any{"bb": map[string]any{"cc": 1}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetStrAnyVal(tt.args.m, tt.args.k, tt.args.v)
			if !reflect.DeepEqual(target, tt.args.m) {
				t.Fatal()
			}
		})
	}
}
