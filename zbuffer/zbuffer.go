package zbuffer

import (
	"sync/atomic"
	"unsafe"

	"github.com/panjf2000/gnet/pkg/pool/ringbuffer"
	"golang.org/x/sys/cpu"
)

const (
	// Number of cells in each chunk. the size is larger than usual CPU cores to reduce hash conflict.
	numChunkCells = 100
	// Number of int64s in each cell. there are 2 pad fields, it should not be too small to avoid waste memory.
	// #nosec G103
	cellCapacity = 6 * unsafe.Sizeof(cpu.CacheLinePad{}) / 8
)

// Int64 is an int64 atomic counter.
type ZBuffer struct {
	cells *[numChunkCells]cell
	index uintptr // index to the n array in each cell
}

// cell is a value container for each cpu core.
type cell struct {
	// We have no way to ensure cache line aligned allocations. so the 2 pads are necessary.
	_ cpu.CacheLinePad
	// The sizeof(cell) shuld be integer multiple of sizeof(cpu.CacheLinePad) to avoid false sharing.
	n [cellCapacity]*ringbuffer.RingBuffer //  ringbuffer.New(1024 * 1000)
	_ cpu.CacheLinePad
}

// chunk is used to saves memory by sharing cells between multiple Int64s
type chunk struct {
	cells     [numChunkCells]cell
	nextIndex uintptr
}

// allocate a new Int64 from the chunk.
func (st *chunk) allocate() ZBuffer {
	for i := atomic.LoadUintptr(&st.nextIndex); i+1 < cellCapacity && atomic.CompareAndSwapUintptr(&st.nextIndex, i, i+1); {
		return ZBuffer{&st.cells, i}
	}
	return ZBuffer{nil, 0}
}

// newChunk creates a new chunk.
func newChunk() *chunk {
	return &chunk{}
}

// the last create chunk. atomic.Pointer is better but it's unavailable until go1.19.
var lastChunk atomic.Value

// MakeInt64 creates a new Int64 object.
// Int64 objects must be created by this function, simply initialized doesn't work.
func MakeInt64() ZBuffer {
	ch, ok := lastChunk.Load().(*chunk)
	if ok {
		ret := ch.allocate()
		if ret.cells != nil {
			return ret
		}
	}
	ch = newChunk()
	ret := ch.allocate() // Must be success because there are no race
	lastChunk.Store(ch)
	return ret
}

//go:linkname getm runtime.getm
func getm() uintptr

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

func ThreadHash() uint {
	m := getm()
	// #nosec G103
	return uint(memhash(unsafe.Pointer(&m), 0, unsafe.Sizeof(m)))
}
