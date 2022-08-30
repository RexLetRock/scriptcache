package data

type Message struct {
	MessageId uint64
	GroupId   uint64
	Data      []byte
	Flags     uint64
	CreatedAt uint64
}
