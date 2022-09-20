package zbuffer

import "github.com/RexLetRock/scriptcache/libs/zcount"

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
		s.makIndex.Inc()
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
