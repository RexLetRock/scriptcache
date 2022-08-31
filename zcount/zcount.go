package zcount

import (
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	cacheLineSize = 64
	maphashSeed   = 42
)

func hash64(x uintptr) uint64 {
	x = ((x >> 33) ^ x) * 0xff51afd7ed558ccd
	x = ((x >> 33) ^ x) * 0xc4ceb9fe1a85ec53
	x = (x >> 33) ^ x
	return uint64(x)
}

func maphash64(s string) uint64 {
	if s == "" {
		return maphashSeed
	}
	strh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return uint64(memhash(unsafe.Pointer(strh.Data), maphashSeed, uintptr(strh.Len)))
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

const cstripes = 64

var ptokenPool sync.Pool

type ptoken struct {
	idx uint32
}

type Counter struct {
	stripes [cstripes]cstripe
}

type cstripe struct {
	c   int64
	pad [cacheLineSize - 8]byte
}

func (c *Counter) IncMax(max int) int64 {
	r := c.Inc()
	if r >= int64(max-1) {
		return c.Reset()
	}
	return r
}

func (c *Counter) Inc() int64 {
	return c.Add(1)
}

func (c *Counter) Dec() int64 {
	return c.Add(-1)
}

func (c *Counter) Add(delta int64) int64 {
	t, ok := ptokenPool.Get().(*ptoken)
	if !ok {
		t = new(ptoken)
		t.idx = uint32(hash64(uintptr(unsafe.Pointer(t))) & (cstripes - 1))
	}
	stripe := &c.stripes[t.idx]
	atomic.AddInt64(&stripe.c, delta)
	ptokenPool.Put(t)
	return c.Value()
}

func (c *Counter) IncZ() {
	c.AddZ(1)
}

func (c *Counter) AddZ(delta int64) {
	t, ok := ptokenPool.Get().(*ptoken)
	if !ok {
		t = new(ptoken)
		t.idx = uint32(hash64(uintptr(unsafe.Pointer(t))) & (cstripes - 1))
	}
	stripe := &c.stripes[t.idx]
	atomic.AddInt64(&stripe.c, delta)
	ptokenPool.Put(t)
}

func (c *Counter) Value() int64 {
	v := int64(0)
	for i := 0; i < cstripes; i++ {
		stripe := &c.stripes[i]
		v += atomic.LoadInt64(&stripe.c)
	}
	return v
}

func (c *Counter) Reset() int64 {
	for i := 0; i < cstripes; i++ {
		stripe := &c.stripes[i]
		atomic.StoreInt64(&stripe.c, 0)
	}
	return int64(0)
}
