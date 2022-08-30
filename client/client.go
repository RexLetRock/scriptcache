package client

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zcount"
	"github.com/gorilla/websocket"

	"scriptcache/colf/message"
)

const NCpu = 1
const NRun = 50_000
const ConnSubMaxCB = 100_000

// var addr = flag.String("url", "localhost:8080", "http service address")

func RunClient() {
	flag.Parse()
	log.SetFlags(0)
	bench()
}

func bench() {
	connSub := CreateConnSub(*addr)
	defer connSub.Close()

	msg := message.Message{MessageId: 6585793445600325728, GroupId: 381870481448962, Data: []byte{0, 0}, Flags: 0, CreatedAt: 1661848717}
	zbench.Run(NRun, NCpu, func(_, _ int) {
		a := connSub.Send(&msg)
		if a == 0 {
			fmt.Printf("ERROR %v \n", a)
		}
	})
}

type ConnSub struct {
	conn     *websocket.Conn
	cbResult [ConnSubMaxCB]chan uint64
	cbCouter zcount.Counter
	channel  chan []byte
}

func CreateConnSub(addr string) *ConnSub {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/message"}
	log.Printf("connecting to %s", u.String())
	p := &ConnSub{}
	p.conn, _, _ = websocket.DefaultDialer.Dial(u.String(), nil)
	p.channel = make(chan []byte, 100_000)

	p.ReceiveLoop()
	p.WriteLoop()

	return p
}

func (c *ConnSub) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *ConnSub) WriteLoop() {
	for r := range c.channel {
		c.conn.WriteMessage(websocket.TextMessage, r)
	}
}

func (c *ConnSub) ReceiveLoop() {
	go func() {
		for {
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if msg != nil {
				m := message.Message{}
				m.Unmarshal(msg[4:])
				mNum := binary.LittleEndian.Uint32(msg[0:4])
				if c.cbResult[mNum] != nil {
					c.cbResult[mNum] <- m.MessageId
				}
			}
		}
	}()
}

func (c *ConnSub) Send(m *message.Message) uint64 {
	b, err := m.MarshalBinary()
	if err != nil {
		return 0
	}

	// Create channel for waiting result
	cbID := c.cbCouter.Inc()
	if cbID > ConnSubMaxCB-10 {
		c.cbCouter.Reset()
	}
	c.cbResult[cbID] = make(chan uint64, 1)
	defer close(c.cbResult[cbID])

	// Write to server
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(cbID))
	bend := append(bs, b...)
	c.channel <- bend

	return <-c.cbResult[cbID]
}
