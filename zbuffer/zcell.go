package zbuffer

import (
	"sync/atomic"
	"time"
)

type ZCell struct {
	data    [cBuffSize]byte
	dataLen int
	wtime   int64
	isrun   int32
}

func (s *ZCell) lock() {
	for {
		if atomic.LoadInt32(&s.isrun) == 0 {
			atomic.StoreInt32(&s.isrun, 1)
			break
		}
		time.Sleep(cTimeLockSleep)
	}
}

func (s *ZCell) unlock() {
	atomic.StoreInt32(&s.isrun, 0)
}

func ZCellCreateMulti(size int) (result []*ZCell) {
	result = make([]*ZCell, size)
	for i := 0; i < size; i++ {
		result[i] = &ZCell{}
	}
	return
}
