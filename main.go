package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/tcp"
)

const connHost = "127.0.0.1:8888"

func main() {
	go tcp.ServerStartViaOptions(connHost)
	time.Sleep(2 * time.Second)
	tcp.ClientStart(connHost)
	select {}
}
