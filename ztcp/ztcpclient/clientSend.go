package ztcpclient

import (
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
	s.startSendloopViaChannel()
}

func (s *TcpClient) startSendloopViaChannel() {
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

				s.conn.Write(s.sendBuffer)
				s.sendBuffer = []byte{}
			}
		}
	}
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
	s.sendBuffer = append(s.slice, []byte(zu.ENDLINE)...)
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
