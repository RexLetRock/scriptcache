package zbufferv2

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
			time.Sleep(100 * time.Millisecond)
			s.Flush()
		}
	}()
}

func (s *ZBuffer) Flush() {
	curTime := ztime.UnixNanoNow()
	for _, pCell := range s.cells {
		if (curTime - atomic.LoadInt64(&pCell.wtime)) > int64(time.Second) {
			if pCell.dataLen > 0 {
				s.handleDo(pCell.data[:pCell.dataLen], true)
				pCell.dataLen = 0
			}
		}
	}
}

func (s *ZBuffer) Write(data []byte) {
	id := getGID()
	pCell := s.cells[id]
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
}

func (s *ZBuffer) handleDo(data []byte, flush bool) {
	if s.handle != nil {
		s.handle(data)
	}
}

func getGID() uint16 {
	return uint16(zgoid.Get()) % cCellSize
}

type ZCell struct {
	data    [cBuffSize]byte
	dataLen int
	wtime   int64
}

func ZCellCreateMulti() (result [cCellSize]*ZCell) {
	for i := 0; i < cCellSize; i++ {
		result[i] = &ZCell{}
	}
	return
}
