package main

import "scriptcache/udp"

func main() {
	// cmd.MultiChannel()
	go udp.ServerStart()
	go udp.ClientStart()

	// go tcpevio.ServerStart()
	// time.Sleep(1 * time.Second)
	// go tcpevio.ClientStart()

	// go tcp.ServerStart()
	// time.Sleep(1 * time.Second)
	// go tcp.ClientStart()

	select {}
}
