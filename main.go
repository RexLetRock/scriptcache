package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zevio"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpclient"
)

const connHost = "127.0.0.1:9000"

func main() {
	// go zgnet.MainGnet()
	go zevio.MainEvio()

	// go ztcpserver.ServerStartViaOptions(connHost)
	time.Sleep(2 * time.Second)
	ztcpclient.ClientStart(connHost)

	select {}
}
