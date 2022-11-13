package zbuffer

import (
	"sync/atomic"
	"time"

	"github.com/RexLetRock/zlib/zgoid"
)

const CFactoryCPUMap = 10_000

type Buffer struct {
	Factories    []*Factory
	FactoriesLen int32

	FactoriesMap   [CFactoryCPUMap]*Factory
	FactoriesIndex Count32
	Handle         func(data []byte)
	Adding         int32
}

func BufferCreate(handle func(data []byte)) *Buffer {
	s := &Buffer{
		Factories:    FactoryCreateMultiple(CMaxCpu, handle),
		Handle:       handle,
		FactoriesLen: CMaxCpu,
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

func (s *Buffer) GetFactory() *Factory {
	id := s.FactoriesIndex.IncMax(CFactoryCPUMap)
	if id >= s.FactoriesLen-1 {
		if atomic.LoadInt32(&s.Adding) == 0 {
			atomic.StoreInt32(&s.Adding, 1)
			s.Factories = append(s.Factories, FactoryCreateMultiple(CMaxCpu, s.Handle)...)
			s.FactoriesLen = int32(len(s.Factories))
			atomic.StoreInt32(&s.Adding, 0)
		} else {
			for {
				if len(s.Factories) >= int(id)+1 {
					break
				}
				time.Sleep(time.Millisecond)
			}
		}
	}

	return s.Factories[id]
}

func (s *Buffer) Write(data []byte) {
	zgoid.Get()
	id := getGID()
	if s.FactoriesMap[id] == nil {
		s.FactoriesMap[id] = s.GetFactory()
	}
	s.FactoriesMap[id].Write(data)
}

func getGID() uint16 {
	return uint16(zgoid.Get())
}
