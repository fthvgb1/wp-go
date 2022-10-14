package safety

import (
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
}
