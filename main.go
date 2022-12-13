package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/ztcp"
)

const cWait = 10 * time.Second

func main() {
	// ztcp.Bench()
	ztcp.BenchMap()
	// zbuffer.ExampleSimple()
	// zbuffer.Bench()
	// zbuffer.ExampleSimple()
	// zbuffer.ExampleSimple()
	// time.Sleep(cWait)
}
