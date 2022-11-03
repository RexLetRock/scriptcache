package ztcpclient

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zcount"

	"github.com/RexLetRock/scriptcache/ztcp/ztcputil"
	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

var count zcount.Counter
var msgf2 = "How are you today ?" // Beware of memleak because buffer

const CMaxResultRetry = 10

var Result ztcputil.ConcurrentMap // [CMaxResultBuffer]*[]byte
var ResultIndex ztcputil.Count32

func init() {
	Result = ztcputil.CMapCreate()
}

type TcpClient struct {
	conn   net.Conn
	chans  chan []byte
	flush  chan []byte
	buffer []byte
	nframe ztcputil.Count32

	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewTcpClient(addr string) *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, zu.ChanSize),
		flush:  make(chan []byte, 1),
		buffer: make([]byte, zu.BuffSize),
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
	o [zu.CRound]*TcpClient
	c ztcputil.Count32
}

func MultiClientCreate(addr string) *MultiClient {
	s := &MultiClient{}
	for i := 0; i < zu.CRound; i++ {
		if s.o[i] = NewTcpClient(addr); s.o[i] == nil {
			panic("cant connect to server")
		}
	}
	return s
}

func (s *MultiClient) Get() *TcpClient {
	return s.o[s.c.Inc()%zu.CRound]
}

func (s *MultiClient) GetViaCpu(cpu int) *TcpClient {
	return s.o[cpu%zu.CRound]
}

func (s *MultiClient) SendMessage(msg string) string {
	return s.Get().SendMessage(msg)
}

func (s *MultiClient) SendMessageViaCpu(msg string, cpu int) string {
	return s.GetViaCpu(cpu).SendMessage(msg)
}

func (s *MultiClient) GetMessage(key string) []byte {
	return s.Get().GetMessage(key)
}

func ClientStart(addr string) {
	var tcpClients = MultiClientCreate(addr)
	time.Sleep(time.Second)

	logrus.Warnf("CLIENT ---msg---> SERVER ---msg---> CLIENT count(msg)")
	logrus.Warnf("Send 50M msg - %v", msgf2)

	groupID := "1"
	// zbench.Run(zu.NRun, zu.NCpu, func(_, j int) {
	// 	tcpClients.SendMessage(MessageNew.Toa() + zu.FRAMESPLIT + groupID)
	// 	// tcpClients.SendMessageViaCpu(MessageNew.Toa()+zu.FRAMESPLIT+groupID, j)
	// })

	// Ticket system
	// logrus.Warn(GetGroupMessageID(tcpClients.GetMessage(tcpClients.SendMessage(MessageNew.Toa() + zu.FRAMESPLIT + groupID))))
	// time.Sleep(15 * time.Second)

	// Broadcast system
	nBroadcast := 1000
	zbench.Run(nBroadcast, 12, func(i, thread int) {
		tcpClients.SendMessage(MessageBroadcast.Toa() + zu.FRAMESPLIT + groupID)
	})

	time.Sleep(2 * time.Second)
	logrus.Warnf("Client broadcast precals %v msg \n", zu.Commaize(int64(nBroadcast)*(ztcputil.CRound+1)))
	logrus.Warnf("Client receive and count %v msg \n", zu.Commaize(count.Value()))

	time.Sleep(20 * time.Second)
}

func GetGroupMessageID(msg []byte) string {
	data := strings.Split(string(msg), zu.FRAMESPLIT)
	if len(data) >= 2 {
		return data[1]
	}
	return ""
}
