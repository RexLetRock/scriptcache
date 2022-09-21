package zbuffer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/zlib/ztime"
	"github.com/sirupsen/logrus"
)

const (
	CPremakeCell     = 2
	CMaxCpu          = 1000 // Goroutine hash
	CMaxChannSize    = 2048 // 1024 * 1kb
	CMaxBuffSize     = 1024 // 1kb
	CMaxOldCellChann = 12
)

type ZBuffer struct {
	Cells       [CMaxCpu]*ZCell
	Factory     [CMaxCpu]*ZCellFactory
	Channs      [CMaxOldCellChann]chan *ZCell
	ChannsIndex zcount.Counter

	cellstime ConcurrentMap
	celldup   ConcurrentMap
	gtime     ztime.Fastime
}

func ZBufferCreate() *ZBuffer {
	var factory [CMaxCpu]*ZCellFactory
	for i := 0; i < CMaxCpu; i++ {
		factory[i] = ZCellFactoryCreate()
	}

	var channs [CMaxOldCellChann]chan *ZCell
	for i := 0; i < CMaxOldCellChann; i++ {
		channs[i] = make(chan *ZCell, 1024)
	}

	s := &ZBuffer{
		Factory:   factory,
		Channs:    channs,
		celldup:   CMapCreate(),
		cellstime: CMapCreate(),
		gtime:     ztime.New(),
	}

	go s.startCellChecktimeLoop()
	go s.startCellChannelLoop()
	go func() {
		for {
			time.Sleep(2 * time.Second) // for i := 0; i < len(s.Factory); i++ { logrus.Warnf("%+v \n", s.Factory[i].Cells) }
			logrus.Warnf("TotalData %v \n", cellDataCount.Value())
		}
	}()
	return s
}

func (s *ZBuffer) startCellChecktimeLoop() {
	for {
		time.Sleep(time.Millisecond)
		curTime := ztime.UnixNanoNow()
		for v := range s.cellstime.IterBuffered() {
			if curTime-v.Val.(int64) >= 5_000_000 {
				key, _ := strconv.Atoi(v.Key)
				s.FlushByTime(s.Cells[key])
			}
		}
	}
}

func (s *ZBuffer) startCellChannelLoop() {
	for _, chann := range s.Channs {
		go func(chann chan *ZCell) {
			for cell := range chann {
				cell.CountData()
				cell = nil
			}
		}(chann)
	}
}

func (s *ZBuffer) ChannGet() chan *ZCell {
	return s.Channs[s.ChannsIndex.Inc()%CMaxOldCellChann]
}

func (s *ZBuffer) FlushByTime(pCell *ZCell) *ZCell {
	if pCell != nil && pCell.dataCount > 0 {
		s.ChannGet() <- pCell
		pCell = s.getCellNewViaID(pCell.hash, pCell.name)
		pCell.dataCount = 0
	}
	return pCell
}

func (s *ZBuffer) Write(data []byte) {
	pData := s.getCell()
	lenData := len(data)
	newDataCount := pData.dataCount + lenData

	// buffer overflow
	if newDataCount >= CMaxBuffSize {
		tmp := pData.data[:pData.dataCount]
		pData.dataCount = 0
		newDataCount = lenData
		select {
		case pData.chann <- tmp:
		default:
			select {
			case s.ChannGet() <- pData:
			default:
				pData = s.getCellNew()
				// pData.chann <- tmp
			}
		}
	}

	copy(pData.data[pData.dataCount:newDataCount], data)
	pData.dataCount = newDataCount
}

func (s *ZBuffer) getCell() *ZCell {
	id, gid := s.getGID()
	if s.Cells[id] == nil {
		s.cellstime.Set(fmt.Sprintf("%v", id), ztime.UnixNanoNow())
		s.Cells[id] = s.Factory[id].GetCell(gid, id)
	}
	return s.Cells[id]
}

func (s *ZBuffer) getCellNew() *ZCell {
	id, gid := s.getGID()
	s.cellstime.Set(fmt.Sprintf("%v", id), ztime.UnixNanoNow())
	s.Cells[id] = s.Factory[id].GetCell(gid, id)
	return s.Cells[id]
}

func (s *ZBuffer) getCellNewViaID(id uint16, gid uint16) *ZCell {
	s.cellstime.Set(fmt.Sprintf("%v", id), ztime.UnixNanoNow())
	s.Cells[id] = s.Factory[id].GetCell(gid, id)
	return s.Cells[id]
}
