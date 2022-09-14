package ztcpclient

import (
	"bufio"

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
				s.handleMsg(msg)
			}
		}
	}()

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
}
