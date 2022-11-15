package zbuffer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/RexLetRock/zlib/zgoid"
	"github.com/sirupsen/logrus"
)

const c1024 = 1024            // For fashion 1024
const cBuffSize = 100 * c1024 // Size of buffer
const cCellSize = 10000       // Number of cell for cpu use

const cTimeLockSleep = 100 * time.Millisecond // Time to sleep before recheck
const cTimeToFlush = 100 * time.Millisecond   // Time to flush when there is not new cell in long time
const cTimeToFlushExit = 1000 * time.Millisecond

var warn = logrus.Warn
var warnf = logrus.Warnf
var skip = func() {}

// Get ID of goroutine
func getGID() int64 {
	return zgoid.Get() % cCellSize
}

// Count32 fast count - fast get
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

// Format number to string
func Commaize(n int64) string {
	s1, s2 := fmt.Sprintf("%d", n), ""
	for i, j := len(s1)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j%3 == 0 && j != 0 {
			s2 = "," + s2
		}
		s2 = string(s1[i]) + s2
	}
	return s2
}
