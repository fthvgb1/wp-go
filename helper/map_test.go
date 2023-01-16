package helper

import (
	"reflect"
	"testing"
)

type Addr struct {
	PostalCode int
	Country    string
}
type Me struct {
	Name    string
	Age     int
	Admin   bool
	Hobbies []string
	Address Addr
	Null    any
}

func TestMapToStruct(t *testing.T) {
	type args struct {
		m map[string]any
	}

	type testCase[T any] struct {
		name    string
		args    args
		wantR   T
		wantErr bool
	}
	tests := []testCase[Me]{
		{
			name: "t1",
			args: args{
				m: map[string]any{
					"name":    "noknow",
					"Age":     2,
					"Admin":   true,
					"Hobbies": []string{"IT", "Travel"},
					"Address": map[string]any{
						"PostalCode": 1111,
						"Country":    "Japan",
					},
					"Null": nil,
				},
			},
			wantR: Me{
				Name:    "noknow",
				Age:     2,
				Admin:   true,
				Hobbies: []string{"IT", "Travel"},
				Address: Addr{
					PostalCode: 1111,
					Country:    "Japan",
				},
				Null: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := MapToStruct[Me](tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("MapToStruct() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	type args[T any] struct {
		s T
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		wantR   map[string]any
		wantErr bool
	}
	tests := []testCase[Me]{
		{
			name: "t1",
			args: args[Me]{
				s: Me{
					Name:    "noknow",
					Age:     2,
					Admin:   true,
					Hobbies: []string{"IT", "Travel"},
					Address: Addr{
						PostalCode: 1111,
						Country:    "Japan",
					},
					Null: nil,
				},
			},
			wantR: map[string]any{
				"Name":    "noknow",
				"Age":     2,
				"Admin":   true,
				"Hobbies": []string{"IT", "Travel"},
				"Address": map[string]any{
					"PostalCode": 1111,
					"Country":    "Japan",
				},
				"Null": nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := StructToMap[Me](tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("StructToMap() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestMapToSlice(t *testing.T) {
	type args[K comparable, V any, T any] struct {
		m  map[K]V
		fn func(K, V) (T, bool)
	}
	type testCase[K comparable, V any, T any] struct {
		name  string
		args  args[K, V, T]
		wantR []T
	}
	tests := []testCase[string, int, int]{
		{
			name: "t1",
			args: args[string, int, int]{
				m: map[string]int{
					"0": 0,
					"1": 1,
					"2": 2,
					"3": 3,
				},
				fn: func(k string, v int) (int, bool) {
					if v > 2 {
						return v, true
					}
					return 0, false
				},
			},
			wantR: []int{3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := MapToSlice(tt.args.m, tt.args.fn); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("MapToSlice() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}