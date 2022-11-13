package zbufferv3

import (
	"sync/atomic"
	"time"

	"github.com/RexLetRock/zlib/zgoid"
	"github.com/RexLetRock/zlib/ztime"
	"github.com/sirupsen/logrus"
)

const c1024 = 1024
const cBuffSize = 100 * c1024 // Size of buffer
const cCellSize = 100         // Number of cell for cpu use

const cTimeLockSleep = 10 * time.Millisecond
const cTimeToFlush = 100 * time.Millisecond

var log = logrus.Warn
var logf = logrus.Warnf

type ZBuffer struct {
	cells  [cCellSize]*ZCell // Cells that store data
	handle func(data []byte) // Function use to handle data pass through buffer
}

func ZBufferCreate(handle func(data []byte)) *ZBuffer {
	s := &ZBuffer{
		handle: handle,
		cells:  ZCellCreateMulti(),
	}

	s.FlushLoopStart()
	return s
}

func (s *ZBuffer) FlushLoopStart() {
	go func() {
		for {
			time.Sleep(cTimeLockSleep)
			s.Flush()
		}
	}()
}

func (s *ZBuffer) Flush() {
	curTime := ztime.UnixNanoNow()
	for _, pCell := range s.cells {
		go func(pCell *ZCell) {
			lastTime := atomic.LoadInt64(&pCell.wtime)
			if lastTime > 0 && (curTime-lastTime) > int64(cTimeToFlush) {
				pCell.lock()
				if pCell.dataLen > 0 {
					s.handleDo(pCell.data[:pCell.dataLen], true)
					pCell.dataLen = 0
				}
				pCell.unlock()
			}
		}(pCell)
	}
}

func (s *ZBuffer) Write(data []byte) {
	id := getGID()
	pCell := s.cells[id]
	pCell.lock()
	atomic.StoreInt64(&pCell.wtime, ztime.UnixNanoNow())

	dataLen := len(data)
	newLen := pCell.dataLen + dataLen
	if newLen >= cBuffSize {
		s.handleDo(pCell.data[:pCell.dataLen], false)
		newLen = dataLen
		pCell.dataLen = 0
	}
	copy(pCell.data[pCell.dataLen:newLen], data)
	pCell.dataLen += dataLen
	pCell.unlock()
}

func (s *ZBuffer) handleDo(data []byte, flush bool) {
	if s.handle != nil {
		s.handle(data)
	}
}

func getGID() uint16 {
	return uint16(zgoid.Get()) % cCellSize
}

// ZCell hold data
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

func ZCellCreateMulti() (result [cCellSize]*ZCell) {
	for i := 0; i < cCellSize; i++ {
		result[i] = &ZCell{}
	}
	return
}
