package main

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/scriptcache/zbuffer"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/sirupsen/logrus"
)

const Address = "127.0.0.1:9000"
const NCpu = 12

var countNum = 0
var countTotal zcount.Counter

func main() {
	// go zgnet.MainGnet()

	// go zevio.MainEvio(Address)
	// go ztcpserver.ServerStartViaOptions(Address)

	// time.Sleep(2 * time.Second)
	// ztcpclient.ClientStart(Address)

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

	// fmt.Println(counter.Read())
	// fmt.Println(rb.Length())
	// fmt.Println(rb.Free())

	zbuffer := zbuffer.ZBufferCreate()
	logrus.Warnf("==== ZBUFFER ===\n")
	fmt.Printf("ZBuffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))
	zbench.Run(200_000_000, 200, func(i, thread int) {
		zbuffer.Write([]byte("Hello How Are You Today|||"))
	})

	go readChann(zbuffer)
	showGoID(zbuffer)
	time.Sleep(10 * time.Second)
	logrus.Errorf("TOTAL SLICE %v \n", countTotal.Value())
	select {}
}

func showGoID(buffer *zbuffer.ZBuffer) {
	count := 0
	for i := 0; i < len(buffer.Cells); i++ {
		if buffer.Cells[i].Name() != 0 {
			count++ // fmt.Printf("%v-", buffer.Cells[i].Name())
		}
	}
	fmt.Printf("\nTOTAL %v \n\n\n", count)
}

func readChann(buffer *zbuffer.ZBuffer) {
	for {
		select {
		case data, ok := <-buffer.Chann:
			if ok {
				a := strings.Split(string(data), "|||")
				countTotal.Add(int64(len(a)))
			}
			// default: // logrus.Warnf("No value ready, moving on.")
		}
	}
}
