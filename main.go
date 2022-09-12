package main

import "github.com/RexLetRock/scriptcache/zgnet"

const connHost = "127.0.0.1:8888"

func main() {
	// go tcp.ServerStartViaOptions(connHost)
	// time.Sleep(2 * time.Second)
	// tcp.ClientStart(connHost)

	zgnet.MainGnet()
	select {}
}
