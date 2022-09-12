package zgnet

import (
	"log"

	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
)

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	logrus.Warnf("Server %v \n", frame)
	out = frame
	return
}

func MainGnet() {
	echo := new(echoServer)
	log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)))
}
