package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"

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

	chans chan []byte
	flush chan []byte

	buffer []byte
	slice  []byte
	conn   net.Conn
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	p := &ConnHandle{
		chans:  make(chan []byte, cChansSize),
		flush:  make(chan []byte, cChansSize),
		buffer: make([]byte, cChansSize),
		conn:   conn,
	}
	p.readerReq, p.writerReq = io.Pipe()

	// Timetoflush
	ticker := time.NewTicker(100 * time.Microsecond)
	go func() {
		for {
			<-ticker.C
			p.flush <- []byte{1}
		}
	}()

	// Receive request msg flow
	go func() {
		reader := bufio.NewReader(p.readerReq)
		for {
			msg, _ := readWithEnd(reader)
			p.chans <- msg
		}
	}()

	// Handle package flow - write flow
	go func() {
		cSend := 0
		for {
			select {
			case msg := <-p.chans:
				p.slice = append(p.slice, msg...)
				cSend += 1
				if cSend >= cSendSize {
					p.conn.Write(p.slice)
					p.slice = []byte{}
					cSend = 0
				}
			case <-p.flush:
				p.conn.Write(p.slice)
				p.slice = []byte{}
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
