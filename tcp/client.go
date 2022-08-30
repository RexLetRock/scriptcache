package tcp

import (
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 12
const NRun = 1_000_000

func ClientStart() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}

	a := []byte("FUCK THIS SHIT IAM DONE")
	zbench.Run(1_000_000, 20, func(i, thread int) {
		conn.Write(a)
	})
	defer conn.Close()
}
