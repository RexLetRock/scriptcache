package ztcpclient

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zcount"

	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
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

	sendWriter *bufio.Writer
}

func NewTcpClient(addr string) *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, zu.ChansSize),
		flush:  make(chan []byte, zu.ChansSize),
		slice:  []byte{},
		buffer: make([]byte, zu.ChansSize),
	}

	s.conn, _ = net.Dial("tcp", addr)
	s.sendWriter = bufio.NewWriter(s.conn)
	s.reader, s.writer = io.Pipe()
	// s.sendReader, s.sendWriter = io.Pipe()

	go s.startTakeloop()
	go s.startSendloop()
	return s
}

func ClientStart(addr string) {
	var tcpClient [zu.NCpu]*TcpClient
	for i := 0; i < zu.NCpu; i++ {
		tcpClient[i] = NewTcpClient(addr)
	}

	logrus.Warnf("CLIENT ---msg---> SERVER ---msg---> CLIENT count(msg)")
	logrus.Warnf("Send 50M msg - empty - channel")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFake()
	})

	logrus.Warnf("Send 50M msg - %v", msgf2)
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFakeV2()
	})

	logrus.Warnf("Send 5M msg - empty - channel - multi-single thread")
	zbench.Run(zu.NRun/10, zu.NCpu, func(i, thread int) {
		tcpClient[0].SendMessageFake()
	})

	time.Sleep(5 * time.Second)
	logrus.Warnf("Client receive and count %v msg \n", zu.Commaize(count.Value()))

}
