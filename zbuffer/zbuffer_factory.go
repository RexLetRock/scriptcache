package zbuffer

import (
	"sync/atomic"
	"time"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/sirupsen/logrus"
)

const CMaxCellPremake = 1
const CMaxCellCircle = 50 // Number of cell
const CMaxCellDelta = 45  // This is gap guard , processing and reusing data -> Delta = Circle - Premake
const CMaxCellOldChannNum = 1

var countAll zcount.Counter

type Cell struct {
	data      []byte
	dataCount int
	name      uint16
	hash      uint16
	pnum      int32
}

type Factory struct {
	ChannsDelta int32
	ChannsIndex int32
	Channs      chan *Cell
	Cells       [CMaxCellCircle]*Cell
	gIndex      int32
	mIndex      int32
}

func CellFactoryCreate() *Factory {
	s := &Factory{Channs: make(chan *Cell, CMaxCellCircle*100)}
	s.MakeCell(CMaxCellPremake)
	go s.handleOldCellLoop()
	return s
}

func (s *Factory) MakeCell(num int) {
	for i := 0; i < num; i++ {
		index := s.mIndex + 1
		s.mIndex = index
		if s.Cells[Gindex(index)] == nil { // cant recycle old cell
			s.Cells[Gindex(index)] = &Cell{
				data:      make([]byte, CMaxBuffSize),
				dataCount: 0,
				name:      0,
				hash:      uint16(index),
			}
		} else {
			p := s.Cells[Gindex(index)]
			p.dataCount = 0
			p.name = 0
			p.hash = uint16(index)
		}
	}
}

func (s *Factory) RateLimit() {
	if s.mIndex-atomic.LoadInt32(&s.ChannsDelta) >= CMaxCellDelta {
		for {
			time.Sleep(time.Millisecond)
			if s.mIndex-atomic.LoadInt32(&s.ChannsDelta) < CMaxCellDelta {
				return
			}
		}
	}
}

func (s *Factory) GetCell(gid uint16, hash uint16) *Cell {
	s.RateLimit()

	// Old cell handle
	select {
	case s.Channs <- s.Cells[s.Gindex()]:
	default:
		logrus.Errorf("Factory chann overflow")
	}

	if s.mIndex-s.gIndex <= CMaxCellPremake {
		s.MakeCell(CMaxCellPremake)
	}

	// Create new cell and use
	s.gIndex++
	pCell := s.Cells[s.Gindex()]
	pCell.name = gid
	pCell.hash = hash
	pCell.pnum = s.gIndex
	return pCell
}

func Gindex(index int32) int32 {
	return index % CMaxCellCircle
}

func (s *Factory) Gindex() int32 {
	return s.gIndex % CMaxCellCircle
}

func (s *Factory) Mindex() int32 {
	return s.mIndex % CMaxCellCircle
}

func log(args ...interface{}) {
	logrus.Warn(args...)
}
