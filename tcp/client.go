package tcp

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"

	"github.com/RexLetRock/scriptcache/colf/message"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 12
const NRun = 5_000_000
const ENDLINE = "#\t#"

var cSend = 0
var cSendSize = 100_000
var chansSize = 1024

type TcpClient struct {
	conn   net.Conn
	chans  chan []byte
	slice  []byte
	buffer []byte

	reader *io.PipeReader
	writer *io.PipeWriter
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
			_, err := readWithEnd(reader)
			if err != nil {
				return
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
			cSend += 1
			s.slice = append(s.slice, msg...)
			if cSend >= cSendSize {
				s.conn.Write(s.slice)
				s.slice = []byte{}
				cSend = 0
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

	// Write to server
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(0))
	bend := append(bs, b...)
	bend = append(bend, []byte(ENDLINE)...)
	c.chans <- bend
	return 0
}

func (c *TcpClient) SendMessageFake() {
	bend := append([]byte{}, []byte(ENDLINE)...)
	c.chans <- bend
}

func ClientStart(addr string) {
	var tcpClient [NCpu]*TcpClient
	for i := 0; i < NCpu; i++ {
		tcpClient[i] = NewTcpClient(addr)
	}
	decodedByteArray, _ := hex.DecodeString("0883E0E9E8FACCD6BE5B1082D8878BD5EF561A0B48656C6C6F206B6974747920DCC0C98080804028C7D4FF888BF1F9024215DCC0C980808040F3C29080808040EF8180808080409001029A010D31363631393439313331323636")
	msg := message.Message{MessageId: 6592524830872596483, GroupId: 382068771122178, Data: decodedByteArray, Flags: 0, CreatedAt: 1661949156780615}
	zbench.Run(NRun, NCpu, func(i, thread int) {
		// tcpClient[thread].SendMessageFake()
		tcpClient[thread].SendMessage(msg)
	})
}
