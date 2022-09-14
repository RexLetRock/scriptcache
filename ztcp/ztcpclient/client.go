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
	bend := append([]byte("How are you today baby"), []byte(zu.ENDLINE)...)
	c.chans <- bend
}

func ClientStart(addr string) {
	var tcpClient [zu.NCpu]*TcpClient
	for i := 0; i < zu.NCpu; i++ {
		tcpClient[i] = NewTcpClient(addr)
	}

	// decodedByteArray, _ := hex.DecodeString("0883E0E9E8FACCD6BE5B1082D8878BD5EF561A0B48656C6C6F206B6974747920DCC0C98080804028C7D4FF888BF1F9024215DCC0C980808040F3C29080808040EF8180808080409001029A010D31363631393439313331323636")
	// decodedByteArray := []byte{}
	// msg := message.Message{MessageId: 6592524830872596483, GroupId: 382068771122178, Data: decodedByteArray, Flags: 0, CreatedAt: 1661949156780615}
	// tcpClient[thread].SendMessage(msg)
	logrus.Warnf("TEST 30M EMPTY")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFake()
	})

	logrus.Warnf("TEST 30M EMPTY")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFake()
	})

	logrus.Warnf("TEST 30M - How are you today baby")
	zbench.Run(zu.NRun, zu.NCpu, func(i, thread int) {
		tcpClient[thread].SendMessageFakeV2()
	})

	time.Sleep(10 * time.Second)
	logrus.Warnf("Msg count %v \n", count.Value())

}
