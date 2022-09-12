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

	iswrite bool
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	s := &ConnHandle{
		chans:   make(chan []byte, cChansSize),
		flush:   make(chan []byte, cChansSize),
		buffer:  make([]byte, cChansSize),
		conn:    conn,
		iswrite: false,
	}
	s.readerReq, s.writerReq = io.Pipe()

	// Timetoflush
	go func() {
		for {
			time.Sleep(cTimeToFlush)
			if s.iswrite {
				s.flush <- []byte{1}
			}
		}
	}()

	// Receive request msg flow
	go func() {
		reader := bufio.NewReader(s.readerReq)
		for {
			msg, _ := readWithEnd(reader)
			s.chans <- msg
		}
	}()

	// Handle package flow - write flow
	go func() {
		cSend := 0
		for {
			select {
			case msg := <-s.chans:
				s.iswrite = true
				s.slice = append(s.slice, msg...)
				cSend += 1
				if cSend >= cSendSize {
					s.conn.Write(s.slice)
					s.slice = []byte{}
					cSend = 0
				}
			case <-s.flush:
				if len(s.slice) > 0 {
					s.iswrite = false
					s.conn.Write(s.slice)
					s.slice = []byte{}
				}
			}
		}
	}()

	// Response flow
	return s
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
