package tcp

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
			tmpI64 := uint64(tmpI)
			if tmpI64 == 0 {
				tmpI64 = cMaxUInt64
			}
			s.result[mNum] = tmpI64
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
			if result == cMaxUInt64 {
				result = 0
			}
			return result
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func ClientStart() {
	tcpClient := NewTcpClient(Addr)

	hexStr := "0883E0E9E8FACCD6BE5B1082D8878BD5EF561A0B48656C6C6F206B6974747920DCC0C98080804028C7D4FF888BF1F9024215DCC0C980808040F3C29080808040EF8180808080409001029A010D31363631393439313331323636"
	decodedByteArray, _ := hex.DecodeString(hexStr)
	msg := message.Message{MessageId: 6592524830872596483, GroupId: 382068771122178, Data: decodedByteArray, Flags: 0, CreatedAt: 1661949156780615} // fmt.Printf("bytes: %b \n val: %v \n str: %s\n", decodedByteArray, decodedByteArray, decodedByteArray)
	// zbench.Run(1, 1, func(i, thread int) {
	// 	_ = tcpClient.SendMessage(msg)
	// })

	fmt.Println(tcpClient.GetMessageID(tcpClient.SendMessage(msg)))
	fmt.Println(tcpClient.GetMessageID(tcpClient.SendMessage(msg)))
	fmt.Println(tcpClient.GetMessageID(tcpClient.SendMessage(msg)))
}
