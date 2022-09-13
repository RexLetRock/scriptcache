package zevio

import (
	"github.com/tidwall/evio"
)

func MainEvio() {
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	if err := evio.Serve(events, "tcp://127.0.0.1:9000"); err != nil {
		panic(err.Error())
	}
}
