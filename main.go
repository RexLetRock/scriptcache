package main

import (
	"github.com/RexLetRock/scriptcache/zbuffer"
)

const Address = "127.0.0.1:9000"

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

	// go readChann(zbuffer)

	// time.Sleep(5 * time.Second)
	// logrus.Errorf("TOTAL SLICE %v \n", countTotal.Value())

	zbuffer.Bench()
	// select {}
}

// func readChann(buffer *zbuffer.ZBuffer) {
// 	for {
// 		select {
// 		case data, ok := <-buffer.Chann:
// 			if ok {
// 				strdata := string(data)
// 				a := strings.Split(strdata, "|||")
// 				countTotal.Add(int64(len(a)))
// 				// logrus.Warnf("STR %v \n", strdata)
// 			}
// 		default: // logrus.Warnf("No value ready, moving on.")
// 		}
// 	}
// }
