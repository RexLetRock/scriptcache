package main

import (
	"github.com/RexLetRock/zlib/zbench"
	"github.com/smallnest/ringbuffer"
)

const Address = "127.0.0.1:9000"
const NCpu = 1

func main() {
	// go zgnet.MainGnet()

	// go zevio.MainEvio(Address)
	// go ztcpserver.ServerStartViaOptions(Address)

	// time.Sleep(2 * time.Second)
	// ztcpclient.ClientStart(Address)
	rb := ringbuffer.New(1024 * 1000)

	zbench.Run(10_000, NCpu, func(i, thread int) {
		rb.Write([]byte("FUCK"))
	})
}
