package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zevio"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpclient"
)

const Address = "127.0.0.1:9000"
const NCpu = 12

func main() {
	// go zgnet.MainGnet()

	go zevio.MainEvio(Address)
	// go ztcpserver.ServerStartViaOptions(Address)

	time.Sleep(2 * time.Second)
	ztcpclient.ClientStart(Address)

	// rb := ringbuffer.New(1024 * 1000 * 1000)
	// var count ztcputil.Count32
	// zbench.Run(50_000_000, NCpu, func(i, thread int) {
	// 	count.Inc()
	// })

	// zbench.Run(50_000_000, NCpu, func(i, thread int) {
	// 	count.Get()
	// })

	// counter := ztcputil.MakeInt64()
	// zbench.Run(50_000_000, NCpu, func(i, thread int) {
	// 	counter.Inc()
	// })

	// a := uint(1)
	// zbench.Run(10000, NCpu, func(i, thread int) {
	// 	b := ztcputil.ThreadHash() % 100
	// 	if a != b {
	// 		a = b
	// 		fmt.Printf("Hash %v \n", a)
	// 	}
	// })

	// fmt.Println(counter.Read())
	// fmt.Println(rb.Length())
	// fmt.Println(rb.Free())
}
