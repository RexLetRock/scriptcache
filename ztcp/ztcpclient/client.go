package ztcpclient

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zcount"
	"github.com/vmihailenco/msgpack/v5"
)

var count zcount.Counter
var msgf2 = "How are you today ?" // Beware of memleak because buffer

type TcpClient struct {
	conn   net.Conn
	chans  chan []byte
	flush  chan []byte
	slice  []byte
	buffer []byte

	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewTcpClient(addr string) *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, zu.ChansSize),
		flush:  make(chan []byte, zu.ChansSize),
		slice:  []byte{},
		buffer: make([]byte, zu.ChansSize),
	}
	s.conn, _ = net.Dial("tcp", addr)
	s.reader, s.writer = io.Pipe()
	go s.startReadLoop()
	go s.startWriteLoop()
	return s
}

func (s *TcpClient) startReadLoop() {
	go func() {
		reader := bufio.NewReader(s.reader)
		for {
			msg, err := zu.ReadWithEnd(reader)
			if err != nil {
				return
			}

			if msg != nil {
				count.Inc()
			}
		}
	}()

	defer s.conn.Close()
	for {
		n, err := s.conn.Read(s.buffer)
		if err != nil {
			return
		}
		s.writer.Write(s.buffer[:n])
	}
}

func (s *TcpClient) startWriteLoop() {
	go func() {
		for {
			time.Sleep(zu.TimeToFlush)
			s.flush <- []byte{}
		}
	}()

	cSend := 0
	for {
		select {
		case msg := <-s.chans:
			cSend += 1
			s.slice = append(s.slice, msg...)
			if cSend >= zu.SendSize {
				go s.conn.Write(s.slice)
				s.slice = []byte{}
				cSend = 0
			}
		case <-s.flush:
			if len(s.slice) > 0 {
				go s.conn.Write((s.slice))
				s.slice = []byte{}
			}
		}
	}
}

func (s *TcpClient) SendBinary(data []byte) {
	s.chans <- data
}

func (c *TcpClient) SendMessage(m interface{}) uint64 {
	b, _ := msgpack.Marshal(m)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(0))
	bend := append(bs, b...)
	bend = append(bend, []byte(zu.ENDLINE)...)
	c.chans <- bend
	return 0
}

func (c *TcpClient) SendMessageFake() {
	c.chans <- []byte(zu.ENDLINE)
}

func (c *TcpClient) SendMessageFakeV2() {
	bend := append([]byte(msgf2), []byte(zu.ENDLINE)...)
	c.chans <- bend
}

func ClientStart(addr string) {
	var tcpClient [zu.NCpu]*TcpClient
	for i := 0; i < zu.NCpu; i++ {
		tcpClient[i] = NewTcpClient(addr)
	}

	logrus.Warnf("CLIENT ---msg---> SERVER ---msg---> CLIENT count(msg)")
	logrus.Warnf("Send 50M msg - empty")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFake()
	})

	logrus.Warnf("Send 50M msg - empty")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFake()
	})

	logrus.Warnf("Send 50M msg - %v ?", msgf2)
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFakeV2()
	})

	time.Sleep(5 * time.Second)
	logrus.Warnf("Client receive and count %v msg \n", zu.Commaize(count.Value()))

}
