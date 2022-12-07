package taskPools

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPools_Execute(t *testing.T) {

	t.Run("pools", func(t *testing.T) {
		p := NewPools(3)
		for i := 0; i < 10; i++ {
			i := i
			p.Execute(func() {
				time.Sleep(time.Second)
				log.Printf("task[%d] done", i)
			})
		}
		p.Wait()
	})
}

func TestPools_ExecuteWithTimeOut(t *testing.T) {

	t.Run("timeout test", func(t *testing.T) {
		p := NewPools(3)
		for i := 0; i < 10; i++ {
			i := i
			p.ExecuteWithTimeOut(2*time.Second, func() {
				log.Printf("start task[%d]", i)
				tt := time.Second
				if i == 0 {
					tt = 7 * time.Second
				}
				time.Sleep(tt)
				log.Printf("task[%d] done", i)
			}, fmt.Sprintf("task [%d]", i))
		}
		p.Wait()
	})
}
