package zbuffer

import (
	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/sirupsen/logrus"
)

type ZCell struct {
	data      []byte
	dataCount int
	name      uint16

	chann chan []byte
}

type ZCellFactory struct {
	Cells    []*ZCell
	name     uint16
	getIndex zcount.Counter
	makIndex zcount.Counter
}

func ZCellFactoryCreate(name uint16) *ZCellFactory {
	s := &ZCellFactory{name: name}
	s.MakeCell(CPremakeCell)
	return s
}

func (s *ZCellFactory) MakeCell(num int) {
	for i := 0; i < num; i++ {
		val := s.makIndex.Inc()
		logrus.Warnf("Make cell gid %v : %v", s.name, val)
		s.Cells = append(s.Cells, &ZCell{
			name:  s.name,
			data:  make([]byte, CMaxBuffSize),
			chann: make(chan []byte, CMaxChannSize),
		})
	}
}

func (s *ZCellFactory) GetCell(gid uint16) *ZCell {
	indexVal := s.getIndex.Inc()
	pCell := s.Cells[indexVal]
	pCell.name = gid

	if s.makIndex.Value()-indexVal <= CPremakeCell+1 {
		s.MakeCell(CPremakeCell)
	}

	return pCell
}
