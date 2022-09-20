package zbuffer

import (
	"fmt"

	"github.com/RexLetRock/zlib/zgoid"
	"github.com/RexLetRock/zlib/ztime"
	"github.com/sirupsen/logrus"
)

const (
	CMaxCpu       = 1000
	CMaxBuffSize  = 1024 * 100
	CMaxChannSize = 1024

	CTimeDiff = 1_000_000 // Milisec to clean old data
)

type ZBuffer struct {
	Cells     [CMaxCpu]ZCell
	dataCount [CMaxCpu]int
	celldup   ConcurrentMap

	Chann chan []byte
}

type ZCell struct {
	data      []byte
	dataCount int
	name      uint16
	time      int64

	Chann chan []byte
}

func ZBufferCreate() *ZBuffer {
	s := &ZBuffer{
		celldup: CMapCreate(),
		Chann:   make(chan []byte, 1024),
	}

	go s.startPullDataLoop()
	return s
}

func (s *ZBuffer) startPullDataLoop() {
	for {
		s.startPullData()
	}
}

func (s *ZBuffer) startPullData() {
	for i := 0; i < CMaxCpu; i++ {
		pCell := s.Cells[i]
		if pCell.name != 0 {
			select {
			case x, ok := <-pCell.Chann:
				if ok {
					s.Chann <- x
				}
			default: //logrus.Warnf("No value ready, moving on.")
			}
		}
	}
}

func (s *ZBuffer) Write(data []byte) {
	pData, gid := s.getCell()
	pData.time = ztime.UnixNanoNow()
	lenData := len(data)
	newDataCount := pData.dataCount + lenData

	// Buffer overflow
	if newDataCount >= CMaxBuffSize {
		pData.dataCount = 0
		newDataCount = lenData
		tmpData := pData.data[:]
		select {
		case pData.Chann <- tmpData:
		default:
			logrus.Errorf("Channel overflow %v \n", gid)
		}
	}

	copy(pData.data[pData.dataCount:newDataCount], data)
	pData.dataCount = newDataCount
}

func (s *ZBuffer) Read() {
}

func (s *ZBuffer) Show() {
	for i, v := range s.Cells {
		if s.dataCount[i] > 0 {
			logrus.Warnf("INFO %v", string(v.data[:s.dataCount[i]]))
		}
	}
}

func (s *ZBuffer) getGID() (id uint16, gid uint16) {
	gid = uint16(zgoid.Get())
	if gid >= CMaxCpu {
		id = gid % CMaxCpu
	} else {
		id = gid
	}
	return
}

func (s *ZBuffer) getCell() (rCell *ZCell, gid uint16) {
	id, gid := s.getGID()
	pCell := &s.Cells[id]
	rCell = pCell

	// Not init yes
	time := ztime.UnixNanoNow()
	if pCell.name == 0 || time-pCell.time > CTimeDiff {
		pCell.time = time
		pCell.name = gid
		pCell.data = make([]byte, CMaxBuffSize)
		pCell.Chann = make(chan []byte, CMaxChannSize)
	} else if pCell.name != gid {
		rCell = s.getCellDup(gid)
	}

	return
}

func (s *ZBuffer) getCellDup(gid uint16) (rCell *ZCell) {
	key := fmt.Sprintf("%v", gid)
	pCell, _ := s.celldup.Get(key)
	if pCell == nil {
		rCell = &ZCell{
			name:  gid,
			data:  make([]byte, CMaxBuffSize),
			Chann: make(chan []byte, CMaxChannSize),
		}
		s.celldup.Set(key, rCell)
	} else {
		rCell = pCell.(*ZCell)
	}
	return
}

func (s *ZCell) Name() uint16 {
	return s.name
}
