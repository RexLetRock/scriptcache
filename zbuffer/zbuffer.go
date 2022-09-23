package zbuffer

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zgoid"
	"github.com/sirupsen/logrus"
)

const (
	CMaxCpu      = 1000      // Goroutine hash
	CMaxBuffSize = 1024 * 10 // 20 * 1kb
	CTimeToFlush = 3000_000
)

type Buffer struct {
	Flush   ConcurrentMap
	Chann   chan uint16
	Cells   []*Cell
	Factory []*Factory
}

func BufferCreate() *Buffer {
	s := &Buffer{
		Flush:   ConcurrentMapCreate(),
		Chann:   make(chan uint16, CMaxCpu),
		Cells:   make([]*Cell, CMaxCpu),
		Factory: FactoryCreateArray(CMaxCpu),
	}

	return s
}

func (s *Buffer) Write(data []byte) {
	cell, id, gid := s.getCell()
	lenData := len(data)
	newDataCount := cell.dataCount + lenData

	if newDataCount >= CMaxBuffSize {
		newDataCount = lenData
		cell = s.getCellNewID(id, gid)
	}

	copy(cell.data[cell.dataCount:newDataCount], data)
	cell.dataCount = newDataCount
}

func (s *Buffer) getCell() (*Cell, uint16, uint16) {
	id, gid := s.getGID()
	if s.Cells[id] == nil {
		s.Cells[id] = s.Factory[id].GetCell(gid, id)
	}
	return s.Cells[id], id, gid
}

func (s *Buffer) getCellNew() *Cell {
	id, gid := s.getGID()
	return s.getCellNewID(gid, id)
}

func (s *Buffer) getCellNewID(id uint16, gid uint16) *Cell {
	s.Cells[id] = s.Factory[id].GetCell(gid, id)
	return s.Cells[id]
}

func (s *Buffer) getGID() (id uint16, gid uint16) {
	gid = uint16(zgoid.Get())
	id = idFromGid(gid)
	return
}

func idFromGid(gid uint16) (id uint16) {
	id = gid
	if gid >= CMaxCpu {
		id = gid % CMaxCpu
	}
	return
}

func Bench() {
	zbuffer := BufferCreate()
	logrus.Warnf("==== ZBUFFER ===\n")
	fmt.Printf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(50_000_000, 12, func(i, thread int) { // zbuffer.Write([]byte("Hello How Are You Today|||"))
		zbuffer.Write([]byte("Hello How Are You Today|||"))
	})

	for {
		time.Sleep(1 * time.Second)
		logrus.Warnf("Countall %v \n", countAll.Value()+400)
	}
}
