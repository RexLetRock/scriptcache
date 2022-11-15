package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zbuffer"
)

const cWait = 10 * time.Second

func main() {
	// ztcp.Bench()
	// zbuffer.ExampleSimple()
	zbuffer.Bench()
	zbuffer.ExampleSimple()
	zbuffer.ExampleSimple()
	time.Sleep(cWait)
}
