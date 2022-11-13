package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zbuffer"
)

const cWait = 10 * time.Second

func main() {
	// ztcp.Bench()
	zbuffer.Bench()
	zbuffer.ExampleSimple()
	time.Sleep(cWait)
}
