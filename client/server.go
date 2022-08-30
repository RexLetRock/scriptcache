package client

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"scriptcache/colf/message"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var channel = make(chan []byte, 1000)

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	go WriteLoop(c)
	go WriteLoop(c)
	go WriteLoop(c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		go func() {
			if msg != nil {

				// DECODE
				m := message.Message{}
				m.Unmarshal(msg[4:])

				// ENCODE
				b, _ := m.MarshalBinary()

				// SEND
				bend := append(msg[0:4], b...)
				channel <- bend
			}
		}()
	}
}

func WriteLoop(c *websocket.Conn) {
	for {
		bend := <-channel
		c.WriteMessage(1, bend)
	}
}

func DoServer() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/message", echo)
	go http.ListenAndServe(*addr, nil)
	time.Sleep(1 * time.Second)
	RunClient()

	bufio.NewScanner(os.Stdin).Scan()
}
