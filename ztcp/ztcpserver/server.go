package ztcpserver

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"

	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

const ThreadPerConn = 5

const cSendSize = 999
const cChanSize = 10000
const cTimeToFlush = time.Millisecond

func ServerStart(host string) {
	log.Printf("Start server at %v\n", host)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Start server err %v \n", err)
	}

	defer listener.Close()
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Printf("Conn accept failed %v \n", err)
		} else {
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

	cntSent int
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	s := &ConnHandle{
		chans:   make(chan []byte, cChanSize),
		flush:   make(chan []byte),
		buffer:  make([]byte, cChanSize),
		conn:    conn,
		cntSent: 0,
	}
	s.readerReq, s.writerReq = io.Pipe()

	// Force Write
	go s.LoopToFlush()

	// Write msg
	go s.LoopToWrite()

	// Receive msg
	go s.LoopToRead()

	return s
}

func (s *ConnHandle) LoopToFlush() {
	for {
		time.Sleep(cTimeToFlush)
		s.flush <- []byte{}
	}
}

func (s *ConnHandle) LoopToWrite() {
	for {
		select {
		case msg := <-s.chans:
			s.slice = append(s.slice, msg...)
			s.cntSent++
			if s.cntSent >= cSendSize {
				if len(s.slice) > 0 {
					go s.conn.Write(s.slice)
					s.slice = []byte{}
					s.cntSent = 0
				}
			}
		case <-s.flush:
			if len(s.slice) > 0 {
				go s.conn.Write(s.slice)
				s.slice = []byte{}
				s.cntSent = 0
			}
		}
	}
}

func (s *ConnHandle) LoopToRead() {
	reader := bufio.NewReader(s.readerReq)
	for {
		msg, _ := zu.ReadWithEnd(reader)
		s.chans <- msg
	}
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
