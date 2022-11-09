package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/ztcp"
)

func main() {
	ztcp.Bench()
	// zbuffer.Bench()
	time.Sleep(10 * time.Second)
}
