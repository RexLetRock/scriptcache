package tcp

import (
	"bufio"
	"io"
	"net"
)

const ThreadPerConn = 5
const countSize = 5_000_000

var pCounter = PerformanceCounterCreate(countSize, 10, "SERVER RUN")

func ServerStart() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer listener.Close()
	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	handle := ConnHandleCreate()
	handle.Handle(conn)
}

type ConnHandle struct {
	reader [ThreadPerConn]*io.PipeReader
	writer [ThreadPerConn]*io.PipeWriter
	buffer []byte
}

func ConnHandleCreate() *ConnHandle {
	p := &ConnHandle{
		buffer: make([]byte, 1024*1000),
	}
	for i := 0; i < ThreadPerConn; i++ {
		p.reader[i], p.writer[i] = io.Pipe()
		go func(index int) {
			reader := bufio.NewReader(p.reader[index])
			for {
				reader.ReadBytes('\n')
				pCounter.Step()
			}
		}(i)
	}
	return p
}

func (s *ConnHandle) Handle(conn net.Conn) error {
	defer conn.Close()
	for {
		n, err := conn.Read(s.buffer)
		if err != nil {
			return err
		}
		s.writer[0].Write(s.buffer[:n])
	}
}
