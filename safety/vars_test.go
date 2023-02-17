package safety

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestVar_Load(t *testing.T) {
	type fields struct {
		val string
		p   unsafe.Pointer
	}
	s := ""
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "t1",
			fields: fields{
				val: s,
				p:   unsafe.Pointer(&s),
			},
			want: "sffs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Var[string]{
				val: tt.fields.val,
				p:   tt.fields.p,
			}
			r.Store(tt.want)
			if got := r.Load(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
	r := NewVar("ff")
	fmt.Println(r.Load())
	q := r
	fmt.Println(q.Load())
	q.Store("xx")
	fmt.Println(r.Load(), q.Load())
}

func TestVar_Delete(t *testing.T) {
	{
		v := NewVar("")
		t.Run("string", func(t *testing.T) {
			v.Delete()
			fmt.Println(v.Load())
			v.Store("xx")
			fmt.Println(v.Load())
		})
	}
}

func TestVar_Flush(t *testing.T) {
	{
		v := NewVar("")
		t.Run("string", func(t *testing.T) {
			v.Flush()
			fmt.Println(v.Load())
			v.Store("xx")
			fmt.Println(v.Load())
		})
	}
}
