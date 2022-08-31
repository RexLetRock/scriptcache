package tcp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"

	"scriptcache/colf/message"
	"scriptcache/zcount"

	"github.com/RexLetRock/zlib/zbench"
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

	cbResult [ConnSubMaxCB]chan uint64
	cbCouter zcount.Counter
}

func NewTcpClient() *TcpClient {
	s := &TcpClient{
		chans:  make(chan []byte, chansSize),
		slice:  []byte{},
		buffer: make([]byte, 1024*1000),
	}
	s.conn, _ = net.Dial("tcp", Addr)
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
			if s.cbResult[mNum] != nil {
				tmpI, _ := strconv.Atoi(string(msg[4 : len(msg)-3]))
				s.cbResult[mNum] <- uint64(tmpI)
			}
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

func (c *TcpClient) SendMessage(m message.Message) (uint64, int64) {
	b, err := m.MarshalBinary()
	if err != nil {
		return 0, 0
	}

	// Create channel for waiting result
	cbID := c.cbCouter.Inc()
	if cbID > ConnSubMaxCB-10 {
		c.cbCouter.Reset()
	}
	c.cbResult[cbID] = make(chan uint64, 1)
	defer close(c.cbResult[cbID])

	// Write to server
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(cbID))
	bend := append(bs, b...)
	bend = append(bend, []byte(ENDLINE)...)
	c.chans <- bend
	// return uint64(cbID)
	return <-c.cbResult[cbID], cbID
}

func ClientStart() {
	tcpClient := NewTcpClient()
	msg := message.Message{MessageId: 6585793445600325728, GroupId: 381870481448962, Data: []byte{0, 0}, Flags: 0, CreatedAt: 1661848717}
	zbench.Run(20_000, 2, func(i, thread int) {
		id, cbid := tcpClient.SendMessage(msg)
		if id != uint64(cbid) {
			fmt.Printf("ERR %v %v \n", id, cbid)
		}
	})
}
