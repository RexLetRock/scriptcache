package ztcpserver

import "strconv"

type Command int

var commands = []string{"MessageNew", "MessageEdit", "NessageBroadcast"}
var commandsLen = Command(len(commands))

const (
	MessageNew Command = iota
	MessageEdit
	MessageBroadcast
)

func (s Command) Toa() string {
	return strconv.Itoa(int(s))
}

func (s Command) Name() (name string) {
	if s < commandsLen {
		name = commands[s]
	}
	return
}

func (s Command) String() (name string) {
	if s < commandsLen {
		name = commands[s]
	}
	return
}
