package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/ztcp/ztcpclient"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpserver"
)

const connHost = "127.0.0.1:9000"

func main() {
	// go zgnet.MainGnet()
	// go zevio.MainEvio()

	go ztcpserver.ServerStartViaOptions(connHost)
	time.Sleep(2 * time.Second)
	ztcpclient.ClientStart(connHost)

	select {}
}
