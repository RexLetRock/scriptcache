package main

import (
	"scriptcache/tcpevio"
	"time"
)

func main() {
	// cmd.MultiChannel()
	// go udp.ServerStart()
	// go udp.ClientStart()

	go tcpevio.ServerStart()
	time.Sleep(1 * time.Second)
	go tcpevio.ClientStart()

	select {}
}
