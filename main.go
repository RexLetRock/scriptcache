package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/tcp"
)

const connHost = "127.0.0.1:9000"

func main() {
	// go tcp.ServerStartViaOptions(connHost)
	// go zevio.MainEvio()
	// go zgnet.MainGnet()
	go tcp.ServerStart()
	time.Sleep(3 * time.Second)
	tcp.ClientStart(connHost)
	select {}
}
