package ztcpclientv2

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

func (s *TcpClient) SendMessageEncode(m interface{}) uint64 {
	b, _ := msgpack.Marshal([2]interface{}{s.nframe.Inc(), m})
	b = append(b, zu.ENDBYTE...)
	s.chans <- b
	return 0
}

func (s *TcpClient) SendMessage() string {
	key := strconv.Itoa(int(ResultIndex.Inc()))
	msg := []byte(key + zu.FRAMESPLIT + msgf2 + zu.ENDLINE)
	s.chans <- msg
	return key
}

func (s *TcpClient) GetMessage(key string) []byte {
	retry := 0
	for {
		if result, _ := Result.Get(key); result != nil {
			Result.Remove(key)
			return *(result.(*[]byte))
		} else {
			time.Sleep(time.Millisecond)
			retry++
			if retry >= CMaxResultRetry {
				return nil
			}
		}
	}
}
