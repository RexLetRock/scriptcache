package zbuffer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/RexLetRock/zlib/ztime"
)

var CellPool = sync.Pool{
	New: func() interface{} { return new(ZCell) },
}

type ZBuffer struct {
	cells      [cCellSize]*ZCell // Cells that store data
	cellsFlush [cCellSize]int32

	handle func(data []byte) // Function use to handle data
	chann  chan *ZCell
}

func ZBufferCreate(handle func(data []byte)) *ZBuffer {
	s := &ZBuffer{
		handle: handle,
		chann:  make(chan *ZCell, c1024),
	}
	go s.startBackgroundJob()
	return s
}

func (s *ZBuffer) Write(data []byte) {
	// Get cell and lock/unlock
	pCell := s.getCellViaPool(getGID())
	pCell.lock()
	defer pCell.unlock()

	// Handle full cell
	dataLen := len(data)
	newLen := pCell.dataLen + dataLen
	if newLen >= cBuffSize {
		atomic.StoreInt64(&pCell.wtime, ztime.UnixNanoNow())
		s.Handle(pCell.data[:pCell.dataLen])
		newLen = dataLen
		pCell.dataLen = 0
	}

	// Handle not full cell
	copy(pCell.data[pCell.dataLen:newLen], data)
	pCell.dataLen += dataLen
}

func (s *ZBuffer) getCellViaPool(GID int64) *ZCell {
	p := s.cells[GID]
	if p == nil {
		s.cells[GID] = CellPool.Get().(*ZCell)
		p = s.cells[GID]
		p.name = GID
	}

	if atomic.LoadInt32(&s.cellsFlush[GID]) == 0 {
		atomic.StoreInt32(&s.cellsFlush[GID], 1)
		s.chann <- p
	}

	return p
}

func (s *ZBuffer) startBackgroundJob() {
	for {
		select {
		case pCell := <-s.chann:
			go s.FlushCell(pCell)
		default:
			skip()
		}
	}
}

func (s *ZBuffer) FlushCell(pCell *ZCell) {
	defer func() { atomic.StoreInt32(&s.cellsFlush[pCell.name], 0) }()

	for {
		time.Sleep(cTimeLockSleep)
		curTime := ztime.UnixNanoNow()
		lastTime := atomic.LoadInt64(&pCell.wtime)
		if (curTime - lastTime) > int64(cTimeToFlush) {
			pCell.lock()
			if pCell.dataLen > 0 {
				s.Handle(pCell.data[:pCell.dataLen])
				pCell.dataLen = 0
				atomic.StoreInt64(&pCell.wtime, ztime.UnixNanoNow())
			}
			pCell.unlock()
		}

		// Expire flush cell not use anymore
		if (curTime - lastTime) > int64(cTimeToFlushExit) {
			break
		}
	}
}

func (s *ZBuffer) Handle(data []byte) {
	if s.handle != nil && len(data) > 0 {
		s.handle(data)
	}
}
