package zbuffer

import (
	"github.com/RexLetRock/zlib/zgoid"
)

type Buffer struct {
	Factories []*Factory
}

func BufferCreate() *Buffer {
	s := &Buffer{
		Factories: FactoryCreateMultiple(CMaxCpu),
	}

	return s
}

func FactoryCreateMultiple(num int) []*Factory {
	var factories = make([]*Factory, num)
	for i := 0; i < num; i++ {
		factories[i] = FactoryCreate(uint16(i))
	}
	return factories
}

func (s *Buffer) Write(data []byte) {
	zgoid.Get()
	id, _ := getGID()
	s.Factories[id].Write(data)
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
