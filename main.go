package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/tcp"
	"github.com/RexLetRock/scriptcache/zevio"
)

const connHost = "127.0.0.1:9000"

func main() {
	// go tcp.ServerStartViaOptions(connHost)
	// go zevio.MainEvio()
	// go zgnet.MainGnet()
	go zevio.MainEvio()
	time.Sleep(2 * time.Second)
	tcp.ClientStart(connHost)
	select {}
}
