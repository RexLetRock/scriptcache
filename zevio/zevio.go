package zevio

import (
	"github.com/tidwall/evio"
)

func MainEvio(address string) {
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	if err := evio.Serve(events, address); err != nil {
		panic(err.Error())
	}
}
