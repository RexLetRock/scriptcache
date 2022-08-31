package tcpevio

import (
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 20
const NRun = 1_000_000

var conns = [NCpu]net.Conn{}

func ClientStart() {
	for i := 0; i < NCpu; i++ {
		conns[i], _ = net.Dial("tcp", "127.0.0.1:5000")
	}

	a := []byte("a")
	zbench.Run(NRun, NCpu, func(i, thread int) {
		conns[thread].Write(a)
	})
}

func padOrTrim(bb []byte, size int) []byte {
	l := len(bb)
	if l == size {
		return bb
	}
	if l > size {
		return bb[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], bb)
	return tmp
}
