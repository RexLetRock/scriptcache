package udp

import (
	"fmt"
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NRun = 1_000_000
const NCpu = 20

var sockets [NCpu]*net.UDPConn

func ClientStart() {
	for i := 0; i < NCpu; i++ {
		socket, _ := net.DialUDP("udp4", nil, &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 8080,
		})
		defer socket.Close()
		sockets[i] = socket
	}

	fmt.Println("Start Client")
	senddata := []byte("abc")
	zbench.Run(NRun, NCpu, func(i, thread int) {
		sockets[thread].Write(senddata)
	})

	// data := make([]byte, 4096)
	// read, remoteAddr, err := socket.ReadFromUDP(data)
	// if err != nil {
	// 	fmt.Println("read data failed!", err)
	// 	return
	// }
	// fmt.Println(read, remoteAddr)
	// fmt.Printf("%s\n", data)
}
