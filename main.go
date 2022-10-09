package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/ztcp"
	"github.com/RexLetRock/zlib/zgoid"
)

const Address = "127.0.0.1:9000"

const CMaxChannSize = 12

type ChannMulti struct {
	c [CMaxChannSize]chan []byte
	z [CMaxChannSize]int
}

func ChannMultiCreate() *ChannMulti {
	s := &ChannMulti{}
	for i := 0; i < CMaxChannSize; i++ {
		s.c[i] = make(chan []byte, 1000)
		go func(i int) {
			for v := range s.c[i] {
				if len(v) > 0 {
					s.z[i]++
				}
			}
		}(i)
	}
	return s
}

func (s *ChannMulti) Write(data []byte) {
	index := zgoid.Get()
	s.c[index%CMaxChannSize] <- data
}

func (s *ChannMulti) Count() int {
	total := 0
	for _, v := range s.z {
		total += v
	}
	return total
}

func main() {
	ztcp.Bench()
	// zbuffer.Bench()

	// channMulti := ChannMultiCreate()
	// zbench.Run(100_000_000, 12, func(_, j int) {
	// 	channMulti.Write([]byte("How are you|||"))
	// })
	// zbench.Run(100_000_000, 24, func(_, j int) {
	// 	channMulti.Write([]byte("How are you|||"))
	// })
	// zbench.Run(100_000_000, 24, func(_, j int) {
	// 	channMulti.Write([]byte("How are you|||"))
	// })
	time.Sleep(10 * time.Second)
	// logrus.Warn("Total ", channMulti.Count())
}
