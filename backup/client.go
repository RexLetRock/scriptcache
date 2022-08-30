package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/RexLetRock/zlib/zcount"
	"github.com/gorilla/websocket"
)

const NCpu = 20

var addr = flag.String("addr", "localhost:8080", "http service address")

func exitAtErr(err error, str ...string) {
	if err != nil {
		log.Fatal(strings.Join(str, " "), err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	for i := 0; i < NCpu; i++ {
		go bench()
	}
	showResult()

	fmt.Printf("\n\nStop with ctrl + c \n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

var counter zcount.Counter
var counterBunch = int64(1_000)
var timeStart = time.Now().Unix()

func bench() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	exitAtErr(err)
	defer c.Close()

	// READ
	done := make(chan struct{})
	n := int64(0)
	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			} else {
				n += 1
				if n%counterBunch == 0 {
					counter.Add(counterBunch)
				}
			}
		}
	}()

	// WRITE
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	func() {
		for {
			c.WriteMessage(websocket.TextMessage, []byte("adlsfhashdflsajdfla h aldjflajsd"))
		}
	}()
}

func showResult() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				timeNow := time.Now().Unix()
				total := counter.Value()
				fmt.Printf("COUNTER %v Msg/s - TOTAL %v msg \n", total/(timeNow-timeStart), total)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
