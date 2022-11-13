package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zbuffer"
)

const cWait = 1 * time.Second

func main() {
	// ztcp.Bench()
	zbuffer.Bench()
	time.Sleep(cWait)
}
