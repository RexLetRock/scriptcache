package main

import (
	"scriptcache/tcp"
	"time"
)

func main() {
	// cmd.MultiChannel()
	// go udp.ServerStart()
	// go udp.ClientStart()

	go tcp.ServerStart()
	time.Sleep(1 * time.Second)
	go tcp.ClientStart()

	select {}
}
