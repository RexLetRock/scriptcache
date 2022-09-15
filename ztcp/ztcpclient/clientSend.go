package ztcpclient

import (
	"bufio"
	"encoding/binary"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

func (s *TcpClient) startSendloop() {
	go func() {
		for {
			time.Sleep(zu.TimeToFlush)
			s.flush <- []byte{}
		}
	}()
	s.sendThroughBuffer()
}

func (s *TcpClient) sendThroughBuffer() {
	b := bufio.NewReader(s.sendReader)
	send := 0
	if b != nil && send != 0 {

	}
	// for {
	// 	v, _, err := b.()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	send++
	// 	s.sendBuffer = append(s.sendBuffer, v...)
	// 	println("The value is " + string(v))
	// }
}

func (s *TcpClient) SendBinary(data []byte) {
	s.chans <- data
}

func (s *TcpClient) SendMessage(m interface{}) uint64 {
	b, _ := msgpack.Marshal(m)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(0))
	bend := append(bs, b...)
	bend = append(bend, []byte(zu.ENDLINE)...)
	s.chans <- bend
	return 0
}

func (s *TcpClient) SendMessageFakeViaBuffer() {
	s.sendCount += 1
	s.sendBuffer = append(s.sendBuffer, []byte(zu.ENDLINE)...)
	if s.sendCount >= zu.SendSize {
		s.conn.Write(s.sendBuffer)
		s.sendBuffer = []byte{}
		s.sendCount = 0
	}
}

func (s *TcpClient) SendMessageFakeViaBufferV2() {
	s.sendCount += 1
	msg := append([]byte(msgf2), []byte(zu.ENDLINE)...)
	s.sendBuffer = append(s.slice, msg...)
	if s.sendCount >= zu.SendSize {
		s.conn.Write(s.sendBuffer)
		s.sendBuffer = []byte{}
		s.sendCount = 0
	}
}

func (s *TcpClient) SendMessageFake() {
	s.chans <- []byte(zu.ENDLINE)
}

func (s *TcpClient) SendMessageFakeV2() {
	bend := append([]byte(msgf2), []byte(zu.ENDLINE)...)
	s.chans <- bend
}
