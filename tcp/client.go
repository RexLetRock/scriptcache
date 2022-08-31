package tcp

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"time"

	"scriptcache/colf/message"
	"scriptcache/zcount"
)

const NCpu = 12
const NRun = 1_000_000
const Addr = "127.0.0.1:8888"
const ConnSubMaxCB = 100_000
const ENDLINE = "#\t#"

var cSend = 0
var cSendSize = 100_000
var chansSize = 1024 * 1000

var cCounter = PerformanceCounterCreate(10_000, 0, "CLIENT RUN")

type TcpClient struct {
	conn   net.Conn
	chans  chan []byte
	slice  []byte
	buffer []byte

	reader *io.PipeReader
	writer *io.PipeWriter

	result   [ConnSubMaxCB]uint64
	cbCouter zcount.Counter
}

func NewTcpClient(addr string) *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, chansSize),
		slice:  []byte{},
		buffer: make([]byte, 1024*1000),
	}
	s.conn, _ = net.Dial("tcp", addr)
	s.reader, s.writer = io.Pipe()

	// READ RESPONSE AND CALLBACK
	go func() {
		reader := bufio.NewReader(s.reader)
		for {
			msg, err := readWithEnd(reader)
			if err != nil {
				return
			}

			cCounter.Step(true)
			mNum := binary.LittleEndian.Uint32(msg[0:4])
			tmpI, _ := strconv.Atoi(string(msg[4 : len(msg)-3]))
			s.result[mNum] = uint64(tmpI)
		}
	}()

	// RECEIVE LOOP
	go func() {
		for {
			n, err := s.conn.Read(s.buffer)
			if err != nil {
				return
			}
			s.writer.Write(s.buffer[:n])
		}
	}()

	// WRITE LOOP
	go func() {
		for {
			msg := <-s.chans
			if true {
				s.conn.Write(msg)
			} else {
				cSend += 1
				s.slice = append(s.slice, msg...)
				if cSend%cSendSize == 0 {
					s.conn.Write(s.slice)
					s.slice = []byte{}
				}
			}
		}
	}()

	return s
}

func (s *TcpClient) Send(data []byte) {
	s.chans <- data
}

func (c *TcpClient) SendMessage(m message.Message) uint64 {
	b, err := m.MarshalBinary()
	if err != nil {
		return 0
	}

	// Create channel for waiting result
	cbID := c.cbCouter.Inc()
	if cbID > ConnSubMaxCB-10 {
		c.cbCouter.Reset()
	}

	// Write to server
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(cbID))
	bend := append(bs, b...)
	bend = append(bend, []byte(ENDLINE)...)
	c.chans <- bend

	return uint64(cbID)
}

func (c *TcpClient) GetMessageID(cbID uint64) uint64 {
	result := uint64(0)
	if cbID == 0 || cbID > ConnSubMaxCB-10 {
		return result
	}

	for {
		result = c.result[cbID]
		if result != 0 {
			return result
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func ClientStart() {
	tcpClient := NewTcpClient(Addr)
	msg := message.Message{MessageId: 6585793445600325728, GroupId: 381870481448962, Data: []byte{0, 0}, Flags: 0, CreatedAt: 1661848717}
	// zbench.Run(1, 1, func(i, thread int) {
	// 	_ = tcpClient.SendMessage(msg)
	// })

	tcpClient.GetMessageID(tcpClient.SendMessage(msg))
	tcpClient.GetMessageID(tcpClient.SendMessage(msg))
}
