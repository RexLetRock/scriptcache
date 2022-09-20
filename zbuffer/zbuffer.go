package zbuffer

import (
	"github.com/RexLetRock/zlib/ztime"
)

const (
	CPremakeCell  = 2
	CMaxCpu       = 1000 // Goroutine hash
	CMaxChannSize = 1024 // 1024 * 1kb
	CMaxBuffSize  = 1024 // 1kb

	CTimeDiff         = 1_000_000 // Milisec to clean old data
	CTimeDiffWrite    = 1_000
	CTimeStopStart    = 10 // At 0 - 10 - 20 micro
	CTimeStopDuration = 5  // Sleep 5 Micro
	CTimeStopRange    = 5
)

type ZBuffer struct {
	Cells   [CMaxCpu]*ZCell
	Factory [CMaxCpu]*ZCellFactory
	Channs  [CMaxCpu]chan []byte

	celldup ConcurrentMap
	gtime   ztime.Fastime
}

func ZBufferCreate() *ZBuffer {
	var factory [CMaxCpu]*ZCellFactory
	for i := 0; i < CMaxCpu; i++ {
		factory[i] = ZCellFactoryCreate(uint16(i))
	}
	s := &ZBuffer{
		Factory: factory,
		celldup: CMapCreate(),
		gtime:   ztime.New(),
	}
	return s
}

func (s *ZBuffer) Write(data []byte) {
	pData := s.getCell()
	lenData := len(data)
	newDataCount := pData.dataCount + lenData

	// Buffer overflow
	if newDataCount >= CMaxBuffSize {
		tmp := pData.data[:pData.dataCount]
		pData.dataCount = 0
		newDataCount = lenData
		select {
		case pData.chann <- tmp:
		default:
			pData = s.getCellNew()
			pData.chann <- tmp
		}
	}

	copy(pData.data[pData.dataCount:newDataCount], data)
	pData.dataCount = newDataCount
}

func (s *ZBuffer) getCell() *ZCell {
	id, gid := s.getGID()
	if s.Cells[id] == nil {
		s.Cells[id] = s.Factory[id].GetCell(gid)
	}
	return s.Cells[id]
}

func (s *ZBuffer) getCellNew() *ZCell {
	id, gid := s.getGID()
	s.Cells[id] = s.Factory[id].GetCell(gid)
	return s.Cells[id]
}
