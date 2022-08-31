package tcp

import (
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 30
const NRun = 3_000_000

var conns = [NCpu]net.Conn{}

func ClientStart() {
	for i := 0; i < NCpu; i++ {
		conns[i], _ = net.Dial("tcp", "127.0.0.1:8888")
	}

	a := []byte("How are you today :D \n")
	zbench.Run(NRun, NCpu, func(i, thread int) {
		conns[thread].Write(a)
	})
}
