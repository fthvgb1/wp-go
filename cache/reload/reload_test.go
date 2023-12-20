package reload

import (
	"fmt"
	"testing"
)

func TestFlushMapVal(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		c := 0
		v := GetAnyValMapBy("key", 2, struct{}{}, func(a struct{}) (int, bool) {
			c++
			return 33, true
		})
		fmt.Println(v)
		FlushMapVal("key", 2)

		v = GetAnyValMapBy("key", 2, struct{}{}, func(a struct{}) (int, bool) {
			fmt.Println("xxxxx")
			return 33, true
		})
		fmt.Println(v)
		FlushAnyVal("key")
		v = GetAnyValMapBy[int, int, struct{}]("key", 2, struct{}{}, func(a struct{}) (int, bool) {
			fmt.Println("yyyy")
			return 33, true
		})
		fmt.Println(v)
	})
}

func TestGetAnyMapFnBys(t *testing.T) {
	var i int
	t.Run("t1", func(t *testing.T) {
		v := GetAnyMapFnBys[int]("name", func(a int) int {
			i++
			return a + 1
		})
		vv := v(1, 2)
		vvv := v(2, 3)
		fmt.Println(vv, vvv)
		v(1, 2)
		FlushMapVal("name", 2)
		v(2, 3)
		fmt.Println(i)
	})
}
