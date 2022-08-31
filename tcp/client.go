package tcp

import (
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 12
const NRun = 20_000_000
const Addr = "127.0.0.1:8888"

var cSend = 0
var cSendSize = 50_000
var chansSize = 1024 * 1000

type TcpClient struct {
	conns net.Conn
	chans chan ([]byte)
}

func NewTcpClient() *TcpClient {
	p := &TcpClient{
		chans: make(chan []byte, chansSize),
	}
	p.conns, _ = net.Dial("tcp", Addr)

	go func() {
		tmpSlice := []byte{}
		for {
			msg := <-p.chans
			cSend += 1
			tmpSlice = append(tmpSlice, msg...)
			if cSend%cSendSize == 0 {
				p.conns.Write(tmpSlice)
				tmpSlice = []byte{}
			}
		}
	}()

	return p
}

func (s *TcpClient) Send(data []byte) {
	s.chans <- data
}

func ClientStart() {
	tcpClient := NewTcpClient()

	a := []byte("How are you today :D \n")
	zbench.Run(NRun, NCpu, func(i, thread int) {
		tcpClient.Send(a)
	})
}
