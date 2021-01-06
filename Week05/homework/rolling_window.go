package homework

import (
	"sync/atomic"
	"time"
)

// RollingWindow 滑动窗口实现
type RollingWindow struct {
	bucket        []int64
	size          int64
	cellDuration  int64
	index         int64
	totalDuration int64
	beginTime     int64
}

type option func(r *RollingWindow)

func SetSize(size int64) option {
	return func(r *RollingWindow) {
		r.bucket = make([]int64, size+1) //
		r.size = size
	}
}

func SetDuration(d time.Duration) option {
	return func(r *RollingWindow) {
		r.totalDuration = int64(d) / int64(time.Millisecond.Nanoseconds())
	}
}

func NewRollingWindow(options ...option) *RollingWindow {
	r := &RollingWindow{
		bucket:        make([]int64, 11),
		size:          10,
		totalDuration: 5 * time.Second.Milliseconds(),
		beginTime:     time.Now().UnixNano() / time.Millisecond.Nanoseconds(),
	}
	for _, op := range options {
		op(r)
	}
	r.cellDuration = int64(r.totalDuration) / r.size
	return r
}

func (r *RollingWindow) Add() {
	passed := time.Now().UnixNano()/time.Millisecond.Nanoseconds() - r.beginTime
	currentIndex := (passed / r.cellDuration) % int64(len(r.bucket))
	if atomic.CompareAndSwapInt64(&r.index, r.index, currentIndex) {
		// 多使用一个窗口，在此处清空尾部
		atomic.StoreInt64(&r.bucket[(currentIndex+1)%(r.size+1)], 0)
	}
	atomic.AddInt64(&r.bucket[currentIndex], 1)
}

func (r *RollingWindow) Count() int64 {
	index := atomic.LoadInt64(&r.index)
	count := int64(0)
	for i := int64(0); i < r.size; i++ {
		count += atomic.LoadInt64(&r.bucket[index])
		index = (index + r.size) % (r.size + 1)
	}
	return count
}
