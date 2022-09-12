package zgnet

import (
	"log"
	"time"

	"github.com/panjf2000/gnet"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, c.Read()...)
	c.ResetBuffer()

	// Use ants pool to unblock the event-loop.
	go func() {
		time.Sleep(1 * time.Second)
		c.AsyncWrite(data)
	}()

	return
}

func MainGnet() {
	echo := &echoServer{}
	log.Fatal(gnet.Serve(echo.EventServer, "tcp://127.0.0.1:8888", gnet.WithMulticore(true)))
}
