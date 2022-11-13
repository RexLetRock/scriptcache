package zbufferv2

import "sync/atomic"

// fast count - fast get
type Count32 int32

func Count32Create() *Count32 {
	return new(Count32)
}

func (c *Count32) IncMaxInt(i int32) int {
	return int(c.IncMax(i))
}

func (c *Count32) IncMax(i int32) int32 {
	a := atomic.AddInt32((*int32)(c), 1)
	if a < i-1 {
		return a
	} else {
		atomic.StoreInt32((*int32)(c), 0)
		return 0
	}
}

func (c *Count32) Inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *Count32) Get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
