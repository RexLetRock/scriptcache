package ztcpserver

import (
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
const cMsgPartsNum = 3

var gIPData zu.ConcurrentMap

func ServerStart(host string) {
	gIPData = zu.CMapCreate()

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
	gIPData.Set(conn.RemoteAddr().String(), s)

	go s.LoopToFlush()
	go s.LoopToWrite()
	go s.LoopToRead()

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
