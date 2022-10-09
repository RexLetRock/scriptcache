package ztcpclient

import (
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zcount"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

var count zcount.Counter
var msgf2 = "How are you today ?" // Beware of memleak because buffer

const CMaxResultBuffer = 1_000_000

func init() {
}

type TcpClient struct {
	conn   net.Conn
	chans  chan []byte
	flush  chan []byte
	buffer []byte
	nframe ztcputil.Count32
	result [CMaxResultBuffer]*[]byte

	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewTcpClient(addr string, result [CMaxResultBuffer]*[]byte) *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, zu.ChanSize),
		flush:  make(chan []byte, 1),
		buffer: make([]byte, zu.ChanSize),
		result: result,
	}

	var err error
	if s.conn, err = net.Dial("tcp", addr); err != nil {
		return nil
	}

	s.reader, s.writer = io.Pipe()
	go s.startTakeloop()
	go s.startSendloop()
	return s
}

type MultiClient struct {
	o  [zu.CRound]*TcpClient
	c  ztcputil.Count32
	R  [CMaxResultBuffer]*[]byte
	RI ztcputil.Count32
}

func MultiClientCreate(addr string) *MultiClient {
	s := &MultiClient{}
	for i := 0; i < zu.CRound; i++ {
		if s.o[i] = NewTcpClient(addr, s.R); s.o[i] == nil {
			panic("cant connect to server")
		}
	}
	return s
}

func (s *MultiClient) Get() *TcpClient {
	return s.o[s.c.IncMax(zu.CRound)]
}

func ClientStart(addr string) {
	var tcpClients = MultiClientCreate(addr)

	logrus.Warnf("CLIENT ---msg---> SERVER ---msg---> CLIENT count(msg)")
	logrus.Warnf("Send 50M msg - %v", msgf2)

	tmp := []byte{}
	key := tcpClients.Get().SendMessageFakeV3(tcpClients.RI.IncMaxInt(CMaxResultBuffer))
	tcpClients.Get().GetMessageResponse(key)
	logrus.Warn("Result ", string(tmp))

	key = tcpClients.Get().SendMessageFakeV3(tcpClients.RI.IncMaxInt(CMaxResultBuffer))
	tcpClients.Get().GetMessageResponse(key)
	logrus.Warn("Result ", string(tmp))

	key = tcpClients.Get().SendMessageFakeV3(tcpClients.RI.IncMaxInt(CMaxResultBuffer))
	tcpClients.Get().GetMessageResponse(key)
	logrus.Warn("Result ", string(tmp))

	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClients.Get().SendMessageFakeV3(tcpClients.RI.IncMaxInt(CMaxResultBuffer))
	})

	time.Sleep(5 * time.Second)
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClients.Get().GetMessageResponse(i % CMaxResultBuffer)
	})

	logrus.Warnf("Client receive and count %v msg \n", zu.Commaize(count.Value()))
}
