package zbuffer

import (
	"strconv"

	"github.com/RexLetRock/zlib/zgoid"
)

type Buffer struct {
	Factories    []*Factory
	FactoriesMap ConcurrentMap
	Handle       func(data []byte)
}

func BufferCreate(handle func(data []byte)) *Buffer {
	s := &Buffer{
		Factories:    FactoryCreateMultiple(CMaxCpu, handle),
		FactoriesMap: CMapCreate(),
		Handle:       handle,
	}
	return s
}

func FactoryCreateMultiple(num int, handle func(data []byte)) []*Factory {
	var factories = make([]*Factory, num)
	for i := 0; i < num; i++ {
		factories[i] = FactoryCreate(uint16(i), handle)
	}
	return factories
}

func (s *Buffer) Write(data []byte) {
	zgoid.Get()
	id, gid := getGID()
	pFac := s.Factories[id]
	if pFac.name == gid || pFac.name == 0xFF {
		s.Factories[id].SetName(gid)
		s.Factories[id].Write(data)
	} else {
		if iFac, ok := s.FactoriesMap.Get(strconv.Itoa(int(gid))); ok {
			iFac.(*Factory).Write(data)
		} else {
			tmpFac := FactoryCreate(gid, s.Handle)
			tmpFac.Write(data)
			s.FactoriesMap.Set(strconv.Itoa(int(gid)), tmpFac)
		}
	}
}

func getGID() (id uint16, gid uint16) {
	gid = uint16(zgoid.Get())
	if gid >= CMaxCpu {
		id = gid % CMaxCpu
	} else {
		id = gid
	}
	return
}
