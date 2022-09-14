package main

import (
	"time"

	tcp "github.com/RexLetRock/scriptcache/ztcp"
)

const connHost = "127.0.0.1:19999"

func main() {
	go tcp.ServerStartViaOptions(connHost)
	// go zgnet.MainGnet()

	// go zevio.MainEvio()
	time.Sleep(2 * time.Second)
	tcp.ClientStart(connHost)

	select {}
}
