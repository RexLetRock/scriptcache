package ztcpclient

import (
	"bufio"
	"strconv"
	"strings"

	zu "github.com/RexLetRock/scriptcache/ztcp/ztcputil"
)

func (s *TcpClient) startTakeloop() {
	go func() {
		reader := bufio.NewReader(s.reader)
		for {
			msg, err := zu.ReadWithEnd(reader)
			if err != nil {
				return
			}

			if msg != nil {
				go s.handleMsg(msg)
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

func (s *TcpClient) handleMsg(msg []byte) {
	count.Inc()
	msgData := strings.Split(string(msg), "|")
	key, err := strconv.Atoi(msgData[0])
	if err != nil {
		return
	}

	s.result[key] = &msg
}
