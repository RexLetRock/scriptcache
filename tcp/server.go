package tcp

import (
	"log"
	"net"
)

func ServerStart() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	bytes := make([]byte, 1024*10)
	for {
		n, err := conn.Read(bytes)
		if err != nil {
			return
		}
		log.Println("Received : ", string(bytes[:n]))
	}
}
