package client

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"

	"scriptcache/colf/message"
)

const NCpu = 20

var addr = flag.String("url", "localhost:8080", "http service address")

func exitAtErr(err error, str ...string) {
	if err != nil {
		log.Fatal(strings.Join(str, " "), err)
	}
}

func RunClient() {
	flag.Parse()
	log.SetFlags(0)
	bench()
}

func bench() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	exitAtErr(err)
	defer c.Close()

	// RECEIVE
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if msg != nil {
				fmt.Printf("RAW %v \n", msg)
				m := message.Message{}
				_, err = m.Unmarshal(msg)
				fmt.Printf("ERR %v \n", err)
				log.Fatal("MSG : ", m)
			}
		}
	}()

	// SEND
	func() {
		for {
			m := message.Message{MessageId: 6585793445600325728, GroupId: 381870481448962, Data: []byte{0, 0}, Flags: 0, CreatedAt: 1661848717}
			b, err := m.MarshalBinary()
			if err != nil {
				return
			}
			c.WriteMessage(websocket.TextMessage, b)
		}
	}()
}
