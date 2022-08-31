package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"scriptcache/colf/message"
)

const ThreadPerConn = 5
const countSize = 5_00_000

var pCounter = PerformanceCounterCreate(countSize, 0, "SERVER RUN")

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
	handle := ConnHandleCreate(conn)
	handle.Handle()
}

type ConnHandle struct {
	readerReq *io.PipeReader
	writerReq *io.PipeWriter

	buffer []byte
	slice  []byte
	conn   net.Conn
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	p := &ConnHandle{
		buffer: make([]byte, 1024*1000),
		conn:   conn,
	}
	p.readerReq, p.writerReq = io.Pipe()

	// Request flow
	go func() {
		reader := bufio.NewReader(p.readerReq)
		cSend := 0
		for {
			msg, _ := readWithEnd(reader)

			// DECODE
			m := message.Message{}
			m.Unmarshal(msg[4 : len(msg)-3])

			// HANDLE PACKAGE DATA

			// INSERT SQL

			// ENCODE
			_, _ = m.MarshalBinary()

			// SEND
			reMsg := append(msg[0:4], []byte(fmt.Sprintf("%v#\t#", binary.LittleEndian.Uint32(msg[0:4])))...)
			cSend += 1
			if true {
				p.conn.Write(reMsg)
			} else {
				p.slice = append(p.slice, reMsg...)
				if cSend%cSendSize == 0 {
					p.conn.Write(p.slice)
					p.slice = []byte{}
				}
			}
			pCounter.Step(true)
		}
	}()

	// Response flow
	return p
}

func (s *ConnHandle) Handle() error {
	defer s.conn.Close()
	for {
		n, err := s.conn.Read(s.buffer)
		if err != nil {
			return err
		}
		s.writerReq.Write(s.buffer[:n])
	}
}
