package zbuffer

import (
	"strings"
	"sync/atomic"
)

func ChannCreateArray(num int) []chan *Cell {
	chann := make([]chan *Cell, num)
	for i := 0; i < num; i++ {
		chann[i] = make(chan *Cell, CMaxCellCircle)
	}
	return chann
}

func FactoryCreateArray(num int) []*Factory {
	factory := make([]*Factory, num)
	for i := 0; i < num; i++ {
		factory[i] = CellFactoryCreate()
	}
	return factory
}

func (s *Factory) handleOldCellLoop() {
	for v := range s.Channs {
		if v != nil && v.dataCount >= 3 {
			dataStr := string(v.data[:v.dataCount])
			a := strings.Split(dataStr, "|||")
			countAll.Add(int64(len(a) - 1))
			atomic.StoreInt32(&s.ChannsDelta, int32(v.pnum))
		}
	}
}
