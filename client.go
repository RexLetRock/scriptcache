// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func exitAtErr(err error, str ...string) {
	if err != nil {
		log.Fatal(strings.Join(str, " "), err)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	exitAtErr(err)

	defer c.Close()

	// READ
	done := make(chan struct{})
	n := 0
	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			} else {
				n += 1
				if n%10_000 == 0 {
					log.Printf("RUN %v\n", n)
				}
			}
		}
	}()

	// WRITE
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	func() {
		for {
			c.WriteMessage(websocket.TextMessage, []byte("Holla kani youwar"))
		}
	}()

	// select {
	// case <-done:
	// 	return
	// case <-interrupt:
	// 	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 	if err != nil {
	// 		log.Println("write close:", err)
	// 		return
	// 	}
	// 	select {
	// 	case <-done:
	// 	case <-time.After(time.Second):
	// 	}
	// 	return
	// }
}
