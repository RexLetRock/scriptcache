package tcp

import (
	"bufio"
	"io"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

const ThreadPerConn = 5
const countSize = 100_000
const connHost = "0.0.0.0:8888"

var pCounter = PerformanceCounterCreate(countSize, 0, "SERVER RUN")

func ServerStart() {
	listener, err := net.Listen("tcp", connHost)
	if err != nil {
		log.Fatalf("Start server %v \n", err)
	}
	defer listener.Close()
	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConn(conn)
		}
	}
}

func ServerStartViaOptions(host string) {
	log.Printf("Start server at %v\n", host)
	listener, _ := net.Listen("tcp", host)
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
			// m := message.Message{}
			// m.Unmarshal(msg[4 : len(msg)-3])

			// SEND
			// mbin, _ := m.MarshalBinary()
			// reMsg := append(msg[0:4], []byte(fmt.Sprintf("%v#\t#", mbin))...)
			// cSend += 1
			reMsg := msg
			p.slice = append(p.slice, reMsg...)
			cSend += 1
			if cSend >= cSendSize {
				p.conn.Write(p.slice)
				p.slice = []byte{}
				cSend = 0
			}
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
