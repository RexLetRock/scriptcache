package ztcpserver

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

const ThreadPerConn = 5
const countSize = 100_000
const connHost = "0.0.0.0:9000"

var pCounter = zu.PerformanceCounterCreate(countSize, 0, "SERVER RUN")
var counter zcount.Counter

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
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Start server err %v \n", err)
	}

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
	s := &ConnHandle{
		chans:  make(chan []byte, zu.ChansSize),
		flush:  make(chan []byte, zu.ChansSize),
		buffer: make([]byte, zu.ChansSize),
		conn:   conn,
	}
	s.readerReq, s.writerReq = io.Pipe()

	// Timetoflush
	go func() {
		for {
			time.Sleep(zu.TimeToFlush)
			s.flush <- []byte{}
		}
	}()

	// Receive request msg flow
	go func() {
		reader := bufio.NewReader(s.readerReq)
		for {
			msg, _ := zu.ReadWithEnd(reader)
			s.chans <- msg
		}
	}()

	// Handle package flow - write flow
	go func() {
		cSend := 0
		for {
			select {
			case msg := <-s.chans:
				counter.Inc()
				s.slice = append(s.slice, msg...)
				cSend += 1
				if cSend >= zu.SendSize {
					go s.conn.Write(s.slice)
					s.slice = []byte{}
					cSend = 0
				}
			case <-s.flush:
				// logrus.Warnf("COUNTER %v \n", counter.Value())
				if len(s.slice) > 0 {
					go s.conn.Write(s.slice)
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
