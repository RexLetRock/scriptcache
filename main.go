package main

import (
	"encoding/binary"
	"fmt"
	"scriptcache/tcp"
	"time"
)

func main() {
	// cmd.MultiChannel()

	// go udp.ServerStart()
	// go udp.ClientStart()

	// go tcpevio.ServerStart()
	// time.Sleep(1 * time.Second)
	// go tcpevio.ClientStart()

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, 6586444308165587067)
	fmt.Println(bs)

	go tcp.ServerStart()
	time.Sleep(1 * time.Second)
	go tcp.ClientStart()

	select {}
}
