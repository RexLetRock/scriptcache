package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/RexLetRock/zlib/zcount"
)

var counter zcount.Counter

func ServerStart() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer listener.Close()
	time.AfterFunc(10*time.Second, func() { fmt.Printf("RECEIVE %v \n", counter.Value()) })
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) error {
	defer conn.Close()
	bytes := make([]byte, 1024*1000)
	for {
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		counter.Add(int64(n))
	}
}
