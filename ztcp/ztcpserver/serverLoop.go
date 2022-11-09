package ztcpserver

import (
	"bufio"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *ConnHandle) LoopToFlush() {
	for {
		time.Sleep(cTimeToFlush)
		s.flush <- []byte{}
	}
}

func (s *ConnHandle) LoopToWrite() {
	for {
		select {
		case msg := <-s.chans:
			s.slice = append(s.slice, msg...)
			s.cntSent++
			if s.cntSent >= cSendSize {
				if len(s.slice) > 0 {
					go s.conn.Write(s.slice)
					s.slice = []byte{}
					s.cntSent = 0
				}
			}
		case <-s.flush:
			if len(s.slice) > 0 {
				go s.conn.Write(s.slice)
				s.slice = []byte{}
				s.cntSent = 0
			}
		}
	}
}

func (s *ConnHandle) LoopToRead() {
	reader := bufio.NewReader(s.readerReq)
	for {
		msg, err := ReadWithEnd(reader)
		if err != nil {
			logrus.Warnf("Msg %v", err)
			continue
		}

		aMsg := strings.Split(string(msg), "|")
		if len(aMsg) < cMsgPartsNum {
			logrus.Warnf("Msg wrong format")
			continue
		}

		// Happy handle msg
		switch aMsg[1] {
		case MessageBroadcast.Toa():
			items := gIPData.Items()
			for _, v := range items {
				cConn := v.(*ConnHandle)
				cConn.chans <- msg
			}
		default:
		}
		s.chans <- msg
	}
}

func ReadWithEnd(reader *bufio.Reader) ([]byte, error) {
	message, err := reader.ReadBytes('#')
	if err != nil {
		return nil, err
	}

	a1, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a1)
	if a1 != '\t' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}

	a2, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a2)
	if a2 != '#' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}
	return message, nil
}
