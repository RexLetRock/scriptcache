package main

import (
	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/sirupsen/logrus"

	"github.com/alphadose/zenq"
)

const Address = "127.0.0.1:9000"
const NCpu = 12

func main() {
	// go zgnet.MainGnet()

	// go zevio.MainEvio(Address)
	// go ztcpserver.ServerStartViaOptions(Address)

	// time.Sleep(2 * time.Second)
	// ztcpclient.ClientStart(Address)

	var count ztcputil.Count32
	zbench.Run(50_000_000, NCpu, func(i, thread int) {
		count.Inc()
	})

	zbench.Run(50_000_000, NCpu, func(i, thread int) {
		count.Get()
	})

	counter := ztcputil.MakeInt64()
	zbench.Run(50_000_000, NCpu, func(i, thread int) {
		counter.Inc()
	})

	logrus.Warnf("=== Zenq ===")
	var zq [NCpu]*zenq.ZenQ[[]byte]
	for i := 0; i < NCpu; i++ {
		zq[i] = zenq.New[[]byte](100_000)
	}
	zbench.Run(1_000_000, NCpu, func(i, thread int) {
		zq[thread].Write([]byte{})
	})

	counter.Set(0)
	zbench.Run(1_000_000, NCpu, func(i, thread int) {
		if _, e := zq[thread].Read(); e {
			counter.Inc()
		}
	})
	logrus.Warnf("=== Zenq get %v items ===", counter.Read())

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
