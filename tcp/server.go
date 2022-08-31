package tcp

import (
	"log"
	"net"
)

func ServerStart() {
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		log.Println(err)
	}

	log.Println("TCP Server is running on port 8888")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Accepted a new TCP connection.")
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	bytes := make([]byte, 1024*10)
	for {
		n, err := conn.Read(bytes)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Received : ", string(bytes[:n]))
	}
}
