package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"scriptcache/client"
	"scriptcache/colf/message"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func errRun(err error, fn func()) {
	if err != nil {
		fn()
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if msg != nil {

			// DECODE
			m := message.Message{}
			_, err = m.Unmarshal(msg)
			if err != nil {
				break
			}

			// ENCODE
			b, err := m.MarshalBinary()
			if err != nil {
				break
			}

			// SEND
			err = c.WriteMessage(mt, b)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	go http.ListenAndServe(*addr, nil)
	time.Sleep(1 * time.Second)
	client.RunClient()

	fmt.Printf("\n\nStop with ctrl + c \n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
