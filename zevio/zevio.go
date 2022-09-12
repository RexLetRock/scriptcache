package zevio

import (
	"github.com/tidwall/evio"
)

func MainEvio() {
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		// logrus.Warnf("IN %v \n", in)
		out = in
		return
	}
	if err := evio.Serve(events, "tcp://127.0.0.1:8888"); err != nil {
		panic(err.Error())
	}
}
