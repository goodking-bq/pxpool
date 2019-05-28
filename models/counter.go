package models

import "sync/atomic"

// Int64Counter 计数器
type Int64Counter int64

// Inc +1
func (c *Int64Counter) Inc(i int64) {
	atomic.AddInt64((*int64)(c), i)
}

// Dec -1
func (c *Int64Counter) Dec(i int64) {
	atomic.AddInt64((*int64)(c), -i)
}

// Get get
func (c *Int64Counter) Get() int64 {
	return atomic.LoadInt64((*int64)(c))
}
