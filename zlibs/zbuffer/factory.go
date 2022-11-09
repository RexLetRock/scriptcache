package zbuffer

import (
	"sync/atomic"
	"time"

	"github.com/RexLetRock/zlib/ztime"
	"github.com/sirupsen/logrus"
)

type Cell struct {
	data      []byte // Buffer data
	dataCount int    // Data count
}

type Factory struct {
	Cells         [CMaxCellCircle]*Cell // Cells bank
	Cell          *Cell                 // Current cell pointer
	Index         int32                 // Current cell index
	Time          int64
	TimeInProcess int64
	Handle        func(data []byte) // Handle data function

	name uint16
	hash uint16

	// Oldcell handle channel
	ChannsDelta Count32
	Channs      chan *Cell

	start bool // Start backgroud service for first time get Cell
}

func FactoryCreate(index uint16, handle func(data []byte)) *Factory {
	s := &Factory{
		name:   0xFF,
		hash:   index % CMaxCpu,
		Channs: make(chan *Cell, CMaxCellDelta),
		Handle: handle,
	}
	go s.handleOldCellLoop()
	return s
}

func (s *Factory) SetName(name uint16) {
	s.name = name
}

func (s *Factory) Write(data []byte) {
	// Get new cell when empty cell
	dataLen := len(data)
	if s.Cell == nil {
		s.CellGet()
	}

	// Get new cell when overflow
	newDataCount := s.Cell.dataCount + dataLen
	if newDataCount >= CMaxBuffSize {
		s.CellGet()
		newDataCount = dataLen
	}

	// Copy data - this is fastest way to write data
	copy(s.Cell.data[s.Cell.dataCount:newDataCount], data)
	s.Cell.dataCount = newDataCount
	s.updateWriteTime()
}

func (s *Factory) CellGet() *Cell {
	// Cell loop when firstime use this cell
	if !s.start {
		s.start = true
		go s.handleFlushCellLoop()
	}

	// Wait channel handle old cell
	if s.Index-s.ChannsDelta.Get() >= CMaxCellDelta {
		for {
			time.Sleep(time.Millisecond)
			if s.Index-s.ChannsDelta.Get() < CMaxCellDelta {
				break
			}
		}
	}

	// Old cell handle
	select {
	case s.Channs <- s.Cells[s.Index%CMaxCellCircle]:
	default:
		logrus.Errorf("Factory chann overflow")
	}

	// Get new or recycle cell
	s.Index++
	hash := s.Index % CMaxCellCircle
	if s.Cells[hash] == nil {
		s.Cells[hash] = &Cell{
			data: make([]byte, CMaxBuffSize),
		}
	}

	// Reset cell for new use
	pCell := s.Cells[hash]
	pCell.dataCount = 0
	s.Cell = pCell
	return pCell
}

func (s *Factory) handleOldCellLoop() {
	for pCell := range s.Channs {
		if pCell != nil {
			if s.Handle != nil {
				s.Handle(pCell.data[:pCell.dataCount])
			}
			s.ChannsDelta.Inc()
		}
	}
}

// Flush cell after time, for cell that not full yet
func (s *Factory) handleFlushCellLoop() {
	time.Sleep(time.Second)
	for {
		time.Sleep(CTimeFlushSleep)
		lastTime := atomic.LoadInt64(&s.Time)
		if lastTime != 0 && ztime.UnixNanoNow()-lastTime > int64(CTimeDiff) {
			s.CellGet()
			s.Time = 0
		}
	}
}

func (s *Factory) updateWriteTime() {
	time := ztime.UnixNanoNow()
	if time != s.TimeInProcess {
		atomic.SwapInt64(&s.Time, time)
	}
	s.TimeInProcess = time
}
