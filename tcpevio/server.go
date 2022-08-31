package tcpevio

import (
	"github.com/tidwall/evio"
)

func ServerStart() {
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		// fmt.Printf("%s \n", in)
		out = in
		return
	}
	if err := evio.Serve(events, "tcp://localhost:5000"); err != nil {
		panic(err.Error())
	}
}
