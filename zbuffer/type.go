package zbuffer

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

const CMaxCellPremake = 1       // Cell to premake for writing
const CMaxCellCircle = 20       // Number of cell in circle
const CMaxCellDelta = 15        // This is gap guard , processing and reusing data -> Delta = Circle - Premake
const CMaxBuffSize = 1024 * 100 // Buffer size in cell -> Make it faster and consume more memory
const CMaxCpu = 100             // For hashing cell

const CTimeDiff = 100 * time.Millisecond       // Time to flush old cell
const CTimeFlushSleep = 100 * time.Millisecond // Time to check need to flush old cell

var log = logrus.Warn

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
