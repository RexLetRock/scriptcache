package ztcpclientv2

import "strconv"

type Command int

const (
	MessageNew Command = iota
	MessageEdit
	MessageBroadcast
)

func (s Command) Toa() string {
	return strconv.Itoa(int(s))
}

func (s Command) String() string {
	switch s {
	case MessageNew:
		return "MessageNew"
	case MessageEdit:
		return "MessageEdit"
	case MessageBroadcast:
		return "MessageBroadcast"
	}

	return "Unknown"
}
