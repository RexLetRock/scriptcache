package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/RexLetRock/zlib/zcount"
)

const waitTime = 15

var counter zcount.Counter

func ServerStart() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer listener.Close()
	time.AfterFunc(waitTime*time.Second, func() { fmt.Printf("RECEIVE %v \n", counter.Value()) })
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
	tmpData := make([]byte, 1024*1000)

	pipeReader, pipeWriter := io.Pipe()
	reader := bufio.NewReader(pipeReader)
	go func() {
		for {
			reader.ReadBytes('\n')
			counter.Inc()
		}
	}()

	for {
		n, err := conn.Read(tmpData)
		if err != nil {
			return err
		}
		pipeWriter.Write(tmpData[:n])
	}
}
