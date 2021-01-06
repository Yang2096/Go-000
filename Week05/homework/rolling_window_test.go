package homework

import (
	"testing"
	"time"
)

func Test_rw(t *testing.T) {
	rw := NewRollingWindow()

	start := time.Now().UnixNano()
	for i := 0; i < 6e4; i++ {
		if i%1000 == 0 {
			time.Sleep(time.Millisecond * 80)
		}
		rw.Add()
	}
	duration := time.Now().UnixNano() - start
	t.Log(duration / time.Millisecond.Nanoseconds()) // 大于5s

	t.Log(rw.Count()) // 55000 并不均匀
}
