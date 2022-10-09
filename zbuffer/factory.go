package zbuffer

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	"github.com/RexLetRock/zlib/ztime"
	"github.com/sirupsen/logrus"
)

const CMaxCellPremake = 1
const CMaxCellCircle = 20 // Number of cell
const CMaxCellDelta = 19  // This is gap guard , processing and reusing data -> Delta = Circle - Premake
const CMaxBuffSize = 20 * 1024
const CMaxCpu = 1000

const CTimeDiff = 200 * 1000_000     // 10 Millisecond
const CTimeFlushTick = 50 * 1000_000 // 10 Millisecond flush

type Cell struct {
	data      []byte
	dataCount int
	pnum      int32
}

type Factory struct {
	Cells         [CMaxCellCircle]*Cell
	Cell          *Cell
	Start         bool
	Index         int32
	Time          int64
	TimeInProcess int64

	name uint16
	hash uint16

	// Oldcell handle channel
	ChannsDelta ztcputil.Count32
	Channs      chan *Cell
}

func FactoryCreate(index uint16) *Factory {
	s := &Factory{
		name:   index,
		hash:   index % CMaxCpu,
		Channs: make(chan *Cell, CMaxCellDelta),
	}
	go s.handleOldCellLoop()
	return s
}

func (s *Factory) WriteTime() {
	time := ztime.UnixNanoNow()
	if time != s.TimeInProcess {
		atomic.SwapInt64(&s.Time, time)
	}
	s.TimeInProcess = time
}

func (s *Factory) Write(data []byte) {
	dataLen := len(data)
	if s.Cell == nil {
		s.CellGet()
	}

	newDataCount := s.Cell.dataCount + dataLen
	if newDataCount >= CMaxBuffSize {
		s.CellGet()
		newDataCount = dataLen
	}

	copy(s.Cell.data[s.Cell.dataCount:newDataCount], data)
	s.Cell.dataCount = newDataCount
	s.WriteTime()
}

func (s *Factory) CellGet() *Cell {
	// Cell loop when firstime use this cell
	if !s.Start {
		s.Start = true
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
	pCell := s.Cells[hash]
	pCell.dataCount = 0
	pCell.pnum = s.Index
	s.Cell = pCell
	return pCell
}

func (s *Factory) handleOldCellLoop() {
	for pCell := range s.Channs {
		if pCell != nil {
			dataStr := string(pCell.data[:pCell.dataCount])
			a := strings.Split(dataStr, "|||")
			countAll.Add(int64(len(a) - 1))
			s.ChannsDelta.Inc()
		}
	}
}

func (s *Factory) handleFlushCellLoop() {
	for {
		time.Sleep(CTimeFlushTick)
		curTime := ztime.UnixNanoNow()
		lastTime := atomic.LoadInt64(&s.Time)
		if lastTime != 0 && curTime-lastTime > CTimeDiff {
			logf()
			s.CellGet()
			s.Time = 0
		}
	}
}
