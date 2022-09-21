package zbuffer

import (
	"strings"

	"github.com/RexLetRock/scriptcache/libs/zcount"
)

var cellDataCount zcount.Counter

type ZCell struct {
	data      []byte
	dataCount int
	name      uint16
	hash      uint16
	chann     chan []byte
}

type ZCellFactory struct {
	Cells    [100_000]*ZCell
	getIndex zcount.Counter
	makIndex zcount.Counter
}

func ZCellFactoryCreate() *ZCellFactory {
	s := &ZCellFactory{}
	s.MakeCell(CPremakeCell)
	return s
}

func (s *ZCellFactory) MakeCell(num int) {
	for i := 0; i < num; i++ {
		s.Cells[s.makIndex.Inc()] = &ZCell{
			data:  make([]byte, CMaxBuffSize),
			chann: make(chan []byte, CMaxChannSize),
		}
	}
}

func (s *ZCellFactory) GetCell(gid uint16, hash uint16) *ZCell {
	indexVal := s.getIndex.Inc()
	pCell := s.Cells[indexVal]
	pCell.name = gid
	pCell.hash = hash

	if s.makIndex.Value()-indexVal < CPremakeCell+3 {
		s.MakeCell(CPremakeCell)
	}

	return pCell
}

func (s *ZCell) CountData() {
	a := strings.Split(string(s.data[:s.dataCount]), "|||")
	cellDataCount.Add(int64(len(a)))

	for {
		select {
		case data := <-s.chann:
			a := strings.Split(string(data), "|||")
			cellDataCount.Add(int64(len(a)))
		default:
			return // logrus.Warnf("End of data %v \n", s.name)
		}
	}

	// Destroy cell
}
