package ztcpclient

import (
	"strconv"
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
	tmpSlice := []byte{}

	for {
		select {
		case msg := <-s.chans:
			cSend += 1
			tmpSlice = append(tmpSlice, msg...)
			if cSend >= zu.SendSize {
				s.conn.Write(tmpSlice)
				tmpSlice = []byte{}
				cSend = 0
			}
		case <-s.flush:
			if len(tmpSlice) > 0 {
				s.conn.Write(tmpSlice)
				tmpSlice = []byte{}
				cSend = 0
			}
		}
	}
}

func (s *TcpClient) SendBinary(data []byte) {
	s.chans <- data
}

func (s *TcpClient) SendMessage(m interface{}) uint64 {
	b, _ := msgpack.Marshal([2]interface{}{s.nframe.Inc(), m})
	b = append(b, []byte(zu.ENDLINE)...)
	s.chans <- b
	return 0
}

func (s *TcpClient) SendMessageFake() {
	s.chans <- []byte(zu.ENDLINE)
}

func (s *TcpClient) SendMessageFakeV2() {
	bend := append([]byte(msgf2), []byte(zu.ENDLINE)...)
	s.chans <- bend
}

func (s *TcpClient) SendMessageFakeV3(key int) int {
	msg := []byte(strconv.Itoa(int(key)) + "|" + msgf2 + zu.ENDLINE)
	s.chans <- msg
	return key
}

func (s *TcpClient) GetMessageResponse(key int) []byte {
	for {
		if res := s.result[key]; res != nil {
			return *res
		}
	}
}
